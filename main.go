package main

import (
	"github.com/miknikif/vault-auto-unseal/command"
	"os"
)

func init() {}

func main() {
	os.Exit(command.Run(os.Args[1:]))
}
