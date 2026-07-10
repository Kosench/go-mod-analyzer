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
