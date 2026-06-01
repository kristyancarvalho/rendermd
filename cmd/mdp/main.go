package main

import (
	"fmt"
	"os"

	"github.com/kristyancarvalho/mdp/internal/app"
)

var version = "dev"

func main() {
	if err := app.Run(version, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "mdp: %v\n", err)
		os.Exit(1)
	}
}
