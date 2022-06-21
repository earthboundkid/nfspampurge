package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/nfspampurge/purge"
)

func main() {
	exitcode.Exit(purge.CLI(os.Args[1:]))
}
