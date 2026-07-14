// Package main is the entry point for the go-mod-analyzer CLI tool.
//
// It contains minimal logic: wiring dependencies and launching.
// All business logic lives in the internal/ directory.

package main

import (
	"os"

	"github.com/Kosench/go-mod-analyzer/internal/cli"
)

var version = "dev"

func main() {
	root := cli.NewRootCmd(os.Stdout)

	// Добавляем флаг --version поверх стандартных команд
	root.Flags().BoolP("version", "V", false, "показать версию и выйти")

	root.Version = version
	root.SetVersionTemplate(`go-mod-analyzer {{.Version}}` + "\n")

	if err := root.Execute(); err != nil {
		// cobra сам печатает ошибки, поэтому просто выходим с кодом 1.
		os.Exit(1)

	}
}
