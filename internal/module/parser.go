package module

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

type Parser interface {
	Parse(path string) ([]Module, error)
}

type modParser struct{}

func NewParser() Parser {
	return &modParser{}
}

func (p *modParser) Parse(path string) ([]Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	// Передаём реальный путь в modfile.Parse для точных сообщений об ошибках.
	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parse modfile %q: %w", path, err)
	}

	modules := make([]Module, 0, len(f.Require))
	for _, req := range f.Require {
		modules = append(modules, Module{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}

	return modules, nil
}
