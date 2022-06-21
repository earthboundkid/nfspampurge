package purge

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/carlmjohnson/flagx"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/versioninfo"
)

const AppName = "nfspampurge"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}
	if err = app.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)

	app.Logger = log.New(os.Stderr, AppName+" ", log.LstdFlags)
	flagx.BoolFunc(fl, "silent", "suppress logging", func() error {
		app.Logger.SetOutput(io.Discard)
		return nil
	})

	appID := fl.String("app-id", "", "`id` for Netlify app")
	formID := fl.String("form-id", "", "`id` for Netlify form")
	cookie := fl.String("cookie", "", "`_nf-auth` value for Netlify cookie")
	fl.DurationVar(&http.DefaultClient.Timeout, "timeout", 5*time.Second, "timeout for connecting to Netlify")
	age := fl.Duration("age", 5*time.Minute, "minimum age for spam comment to purge")
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `Netlify Spam Purge - %s

Deletes all messages in Netlify's spam box.
Options may be passed as env vars like NFSPAMPURGE_APP_ID.

Get the cookie value by running this in the dev console:

copy(JSON.parse(localStorage.getItem("nf-session")).access_token)

Usage:

	nfspampurge [options]

Options:
`, versioninfo.Short())
		fl.PrintDefaults()
	}
	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagx.ParseEnv(fl, AppName); err != nil {
		return err
	}
	if err := flagx.MustHave(fl, "app-id", "form-id", "cookie"); err != nil {
		return err
	}
	app.cl = requests.
		URL("https://app.netlify.com/").
		Pathf(
			"/access-control/bb-api/api/v1/sites/%s/forms/%s/submissions",
			*appID, *formID,
		).
		Cookie("_nf-auth", *cookie)
	app.before = time.Now().Add(-*age)
	return nil
}

type appEnv struct {
	*log.Logger
	cl     *requests.Builder
	before time.Time
}

func (app *appEnv) Exec() (err error) {
	app.Println("starting")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	entries, err := app.Get(ctx)
	if err != nil {
		return err
	}
	{
		filteredEntries := entries[:0]
		for _, entry := range entries {
			if entry.CreatedAt.Before(app.before) {
				filteredEntries = append(filteredEntries, entry)
			}
		}
		entries = filteredEntries
		app.Printf("%d entries from before %s", len(entries), app.before.Format(time.RFC1123))
	}
	if err = app.Purge(ctx, entries); err != nil {
		return err
	}
	return nil
}

func (app *appEnv) Get(ctx context.Context) (entries NFResponse, err error) {
	n := 1
	for {
		var req NFResponse
		err = app.cl.Clone().
			Param("state", "spam").
			Param("per_page", "100").
			ParamInt("page", n).
			ToJSON(&req).
			Fetch(ctx)
		if err != nil {
			return nil, err
		}
		app.Printf("listing %d entries", len(req))
		if len(req) == 0 {
			return entries, nil
		}
		entries = append(entries, req...)
		n++
	}
}

func (app *appEnv) Purge(ctx context.Context, entries NFResponse) (err error) {
	for _, entry := range entries {
		id := entry.ID
		app.Println("purging", id, "from", entry.CreatedAt.Local().Format(time.RFC1123))
		err = app.cl.Clone().
			Pathf("/access-control/bb-api/api/v1/submissions/%s", id).
			Delete().
			CheckStatus(http.StatusNoContent).
			Fetch(ctx)
		if err != nil {
			return err
		}
	}
	return
}

type NFResponse []struct {
	Number             int                  `json:"number"`
	Title              any                  `json:"title"`
	Email              string               `json:"email"`
	Name               string               `json:"name"`
	FirstName          string               `json:"first_name"`
	LastName           any                  `json:"last_name"`
	Company            any                  `json:"company"`
	Summary            string               `json:"summary"`
	Body               string               `json:"body"`
	Data               Data                 `json:"data"`
	CreatedAt          time.Time            `json:"created_at"`
	HumanFields        HumanFields          `json:"human_fields"`
	OrderedHumanFields []OrderedHumanFields `json:"ordered_human_fields"`
	ID                 string               `json:"id"`
	FormID             string               `json:"form_id"`
	SiteURL            string               `json:"site_url"`
	FormName           string               `json:"form_name"`
}

type Data struct {
	HostPage  string `json:"host_page"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Cc        string `json:"CC"`
	Comment   string `json:"comment"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Referrer  string `json:"referrer"`
}

type HumanFields struct {
	HostURL string `json:"Host URL"`
	From    string `json:"From"`
	Email   string `json:"Email"`
	Subject string `json:"Subject"`
	Cc      string `json:"Cc"`
	Comment string `json:"Comment"`
}

type OrderedHumanFields struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
