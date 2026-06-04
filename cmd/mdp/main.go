package main

import (
	"fmt"
	"os"

	"github.com/kristyancarvalho/mdp/internal/app"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	build := app.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	if err := app.Run(build, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "mdp: %v\n", err)
		os.Exit(1)
	}
}
