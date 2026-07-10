package module

import (
	"path/filepath"
	"testing"
)

// TestParseFile tests the go.mod parser.
//
// It uses a table-driven approach: several scenarios are described
// in a slice of structs and run through the same code.

func TestParseFile(t *testing.T) {
	examplePath := filepath.Join("testdata", "go.mod")

	tests := []struct {
		name      string
		path      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "успешный парсинг реального go.mod",
			path:      examplePath,
			wantCount: 4, // cobra, viper, mod, testify
			wantErr:   false,
		},
		{
			name:      "успешный парсинг по пути к каталогу",
			path:      "testdata", // каталог, а не файл
			wantCount: 4,
			wantErr:   false,
		},
		{
			name:      "несуществующий файл возвращает ошибку",
			path:      "testdata/nonexistent.go.mod",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(got) != tt.wantCount {
				t.Fatalf("ParseFile() вернул %d модулей, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestParseFile_Content(t *testing.T) {
	modules, err := ParseFile(filepath.Join("testdata", "go.mod"))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	// Ожидаем cobra первой зависимостью (порядок = порядок в go.mod).
	cobra := modules[0]
	if cobra.Path != "github.com/spf13/cobra" {
		t.Errorf("Path = %q, want %q", cobra.Path, "github.com/spf13/cobra")
	}
	if cobra.Version != "v1.8.0" {
		t.Errorf("Version = %q, want %q", cobra.Version, "v1.8.0")
	}
	if cobra.Indirect {
		t.Error("cobra должна быть прямой зависимостью, got Indirect=true")
	}

	// testify — последняя, indirect.
	testify := modules[3]
	if !testify.Indirect {
		t.Error("testify должна быть indirect, got Indirect=false")
	}
}

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
			if got := m.Name(); got != tt.want {
				t.Errorf("Name() = %q, want %q", got, tt.want)
			}
		})
	}
}
