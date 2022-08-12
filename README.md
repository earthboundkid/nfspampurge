# Netlify Spam Purge [![GoDoc](https://godoc.org/github.com/carlmjohnson/nfspampurge?status.svg)](https://godoc.org/github.com/carlmjohnson/nfspampurge) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/nfspampurge)](https://goreportcard.com/report/github.com/carlmjohnson/nfspampurge)

Deletes all messages in Netlify's spam box.

Usage:

- Go to Netlify's spam page.
- Remove any non-spam from the spam box.
- Get the cookie value by creating a personal access token at https://app.netlify.com/user/applications#personal-access-tokens
- Run `nfspampurge -cookie COOKIE_VALUE` or set the cookie value to the NFSPAMPURGE_COOKIE environment variable.

To prevent new incoming messages from being deleted before being reviewed, by default messages must be at least 5 minutes old to be purged.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" go install github.com/carlmjohnson/nfspampurge@latest
```

## Screenshots

```
$ nfspampurge -h
Netlify Spam Purge - v0.22.2

Deletes all messages in Netlify's spam box.
Options may be passed as env vars like NFSPAMPURGE_APP_ID.

Get the cookie value by creating a personal access token at https://app.netlify.com/user/applications#personal-access-tokens

Usage:

        nfspampurge [options]

Options:
  -age duration
        minimum age for spam comment to purge (default 5m0s)
  -app-id id
        id for Netlify app
  -cookie token
        personal access token value for Netlify cookie
  -form-id id
        id for Netlify form
  -silent
        suppress logging
  -timeout duration
        timeout for connecting to Netlify (default 5s)

$ nfspampurge
nfspampurge 2022/06/21 13:06:24 starting
nfspampurge 2022/06/21 13:06:25 listing 6 entries
nfspampurge 2022/06/21 13:06:25 listing 0 entries
nfspampurge 2022/06/21 13:06:25 5 entries from before Tue, 21 Jun 2022 12:41:24 EDT
nfspampurge 2022/06/21 13:06:25 purging 62b1f4a3ae850c0ff1cdcaef from Tue, 21 Jun 2022 12:41:07 EDT
nfspampurge 2022/06/21 13:06:25 purging 62b1f1cd0108fd0e4aa5e9b7 from Tue, 21 Jun 2022 12:29:01 EDT
nfspampurge 2022/06/21 13:06:25 purging 62b1ec8f84d3d5091c982a98 from Tue, 21 Jun 2022 12:06:39 EDT
nfspampurge 2022/06/21 13:06:26 purging 62b1eab0ad935f0982b9135a from Tue, 21 Jun 2022 11:58:40 EDT
nfspampurge 2022/06/21 13:06:26 purging 62b1e3b72600930bf32a0692 from Tue, 21 Jun 2022 11:28:55 EDT
```
