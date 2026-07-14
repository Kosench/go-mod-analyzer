package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestNewRootCmd_Help проверяет, что help корректно формируется.
// Это "дымовой тест" — убеждаемся, что cobra-команда вообще собрана правильно.
func TestNewRootCmd_Help(t *testing.T) {
	// bytes.Buffer заменяет stdout — пишем сюда и проверяем содержимое.
	buf := &bytes.Buffer{}
	cmd := NewRootCmd(buf)

	// Эмулируем вызов "go-mod-analyzer --help".
	cmd.SetArgs([]string{"--help"})
	cmd.SetOut(buf)

	// Execute вернёт ошибку (help — это special case, cobra возвращает nil,
	// но флаг help выставляется). Нам важен сам факт выполнения без паники.
	_ = cmd.Execute()
	out := buf.String()

	// Должны увидеть Short-описание и флаг format.
	if !strings.Contains(out, "Анализ зависимостей") {
		t.Errorf("help не содержит Short, got:\n%s", out)
	}
	if !strings.Contains(out, "format") {
		t.Errorf("help не содержит флаг format, got:\n%s", out)
	}
}

// TestRunAnalyze проверяет основной сценарий: передали путь, получили вывод.
func TestRunAnalyze(t *testing.T) {
	// Используем testdata из пакета module. Путь относительный
	// от каталога теста (internal/cli/), поэтому поднимаемся на уровень.
	// .../internal/module/testdata
	buf := &bytes.Buffer{}
	opts := Options{
		Path:   "../module/testdata",
		Format: "text",
	}

	cmd := NewRootCmd(buf)
	err := runAnalyze(cmd, buf, opts)
	if err != nil {
		t.Fatalf("runAnalyze вернул ошибку: %v", err)
	}

	out := buf.String()
	// Ожидаем в выводе модули из example.go.mod
	for _, want := range []string{"cobra", "viper", "testify"} {
		if !strings.Contains(out, want) {
			t.Errorf("вывод не содержит %q, got:\n%s", want, out)
		}
	}
}

// TestRunAnalyze_InvalidFormat проверяет обработку неверного формата.
func TestRunAnalyze_InvalidFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Path:   "../module/testdata",
		Format: "xml", // несуществующий
	}

	cmd := NewRootCmd(buf)
	err := runAnalyze(cmd, buf, opts)

	if err == nil {
		t.Fatal("ожидали ошибку, got nil")
	}
	if !strings.Contains(err.Error(), "неподдерживаемый формат") {
		t.Errorf("ожидали ошибку про формат, got: %v", err)
	}
}

// TestRunAnalyze_NonExistentPath проверяет обработку несуществующего пути.
func TestRunAnalyze_NonExistentPath(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Path:   "/nonexistent/path/xyz",
		Format: "text",
	}

	cmd := NewRootCmd(buf)
	err := runAnalyze(cmd, buf, opts)

	if err == nil {
		t.Fatal("ожидали ошибку, got nil")
	}
}
