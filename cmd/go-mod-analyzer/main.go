// Package main is the entry point for the go-mod-analyzer CLI tool.
//
// It contains minimal logic: wiring dependencies and launching.
// All business logic lives in the internal/ directory.

package main

import (
	"os"

	"github.com/Kosench/go-mod-analyzer/internal/cli"
	"github.com/Kosench/go-mod-analyzer/internal/module"
)

var version = "dev"

func main() {
	parser := module.NewParser()

	root := cli.NewRootCmd(os.Stdout, parser, version)

	root.Flags().BoolP("version", "V", false, "show version and exit")
	root.Version = version
	root.SetVersionTemplate("go-mod-analyzer {{.Version}}\n")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
