// Package main is the entry point for the go-mod-analyzer CLI tool.
//
// It contains minimal logic: wiring dependencies and launching.
// All business logic lives in the internal/ directory.

package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "go-mod-analyzer: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("go-mod-analyzer %s - dependency analysis\n", version)
	fmt.Println("TODO: реализовать парсинг go.mod (Этап 1)")
	return nil
}
