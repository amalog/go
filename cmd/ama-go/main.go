package main

import (
	"os"

	"github.com/amalog/go"
)

func main() {
	ama := amalog.Amalog{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	exitCode := ama.Run(os.Args[1:])
	os.Exit(exitCode)
}
