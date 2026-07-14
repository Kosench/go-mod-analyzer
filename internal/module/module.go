// Package module provides domain types and logic for working with valid.mod.

package module

type Module struct {
	Path     string
	Version  string
	Indirect bool
}

func (m Module) Name() string {
	for i := len(m.Path) - 1; i >= 0; i-- {
		if m.Path[i] == '/' {
			return m.Path[i+1:]
		}
	}
	return m.Path
}
