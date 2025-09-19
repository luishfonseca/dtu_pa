// Package main is the entrypoint for the program analyzer.
//
// See the cmd package for more information about the tool.
package main

import (
	"os"

	"github.com/luishfonseca/dtu_pa/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
