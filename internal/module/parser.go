package module

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func ParseFile(path string) ([]Module, error) {
	if filepath.Base(path) != "go.mod" {
		path = filepath.Join(path, "go.mod")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("parse go.mod: %w", err)
	}

	var modules []Module
	for _, req := range f.Require {
		modules = append(modules, Module{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}

	return modules, nil
}
