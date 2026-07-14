// Package module provides domain types and logic for working with valid.mod.

package module

type Module struct {
	Path    string
	Version string
	// Indirect is true if the dependency is pulled in transitively
	// (not explicitly listed but required by another dependency).
	// In valid.mod it is marked with an // indirect comment.
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
