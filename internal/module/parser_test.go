package module

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		wantModules []Module
		wantErr     bool
		errContains string // подстрока в ошибке для точной проверки
	}{
		{
			name:     "valid valid.mod file",
			filePath: filepath.Join("testdata", "valid.mod"),
			wantModules: []Module{
				{Path: "github.com/spf13/cobra", Version: "v1.8.0", Indirect: false},
				{Path: "github.com/stretchr/testify", Version: "v1.8.4", Indirect: true},
			},
			wantErr: false,
		},
		{
			name:        "invalid valid.mod syntax",
			filePath:    filepath.Join("testdata", "invalid.mod"),
			wantErr:     true,
			errContains: "parse modfile",
		},
		{
			name:        "non-existent file returns error",
			filePath:    filepath.Join("testdata", "does_not_exist.mod"),
			wantErr:     true,
			errContains: "read file",
		},
		{
			// ВАЖНЫЙ КЕЙС: доказываем, что парсер не занимается резолвом путей.
			name:        "directory instead of file returns error",
			filePath:    "testdata", // это папка, а не файл
			wantErr:     true,
			errContains: "read file",
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.filePath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains,
						"ошибка должна содержать контекст операции")
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantModules, got,
				"распарсенные модули должны точно совпадать с ожидаемыми")
		})
	}
}

// TestModuleName проверяет метод Name() отдельно.
func TestModuleName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"github.com/spf13/cobra", "cobra"},
		{"golang.org/x/mod", "mod"},
		{"single", "single"}, // нет слеша — весь путь
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			m := Module{Path: tt.path}
			assert.Equal(t, tt.want, m.Name())
		})
	}
}

// TestModParser_InterfaceCompliance — compile-time проверка,
// что *modParser реализует интерфейс Parser.
func TestModParser_InterfaceCompliance(t *testing.T) {
	var _ Parser = NewParser()
}
