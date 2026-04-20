package main

import (
	"os"

	"github.com/cyberbuff/orkl"
)

func main() {
	os.Exit(orkl.RunCLI(os.Args[1:], os.Stdout, os.Stderr))
}
