# go-cli [![GoDoc](https://godoc.org/github.com/carlmjohnson/nfspampurge?status.svg)](https://godoc.org/github.com/carlmjohnson/nfspampurge) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/nfspampurge)](https://goreportcard.com/report/github.com/carlmjohnson/nfspampurge)

Deletes all messages in Netlify's spam box.

Usage: Go to Netlify's spam page. Remove any non-spam from the spam box. Get the cookie value by entering this in the developer console:

```
copy(JSON.parse(localStorage.getItem("nf-session")).access_token)
```

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" go install github.com/carlmjohnson/nfspampurge@latest
```

## Screenshots

```
$ nfspampurge -h
Netlify Spam Purge - devel

Deletes all messages in Netlify's spam box.
Options may be passed as env vars like NFSPAMPURGE_APP_ID.

Get the cookie value by running this in the dev console:

copy(JSON.parse(localStorage.getItem("nf-session")).access_token)

Usage:

        nfspampurge [options]

Options:
  -app-id id
        id for Netlify app
  -cookie _nf-auth
        _nf-auth value for Netlify cookie
  -form-id id
        id for Netlify form
  -silent
        suppress logging
  -timeout duration
        timeout for connecting to Netlify (default 5s)

$ nfspampurge
nfspampurge 2022/06/21 10:36:37 got 55 entries
nfspampurge 2022/06/21 10:36:37 purging 5de5d0d7cc558333db05a492 from Tue, 03 Dec 2019 03:04:55 UTC
nfspampurge 2022/06/21 10:36:37 purging 5de4123fc25e8d9a36938999 from Sun, 01 Dec 2019 19:19:27 UTC
nfspampurge 2022/06/21 10:36:37 purging 5de35d6eb4bbcd961eefa221 from Sun, 01 Dec 2019 06:27:58 UTC
…
nfspampurge 2022/06/21 10:36:51 got 0 entries
```
