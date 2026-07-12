package report

import (
	"strings"
	"testing"
)

// TestTextRenderer_Render проверяет базовый рендер.
func TestTextRenderer_Render(t *testing.T) {
	r := NewTextRenderer()

	// Готовим Summary с двумя модулями, один с уязвимостью.
	used := true
	s := Summary{
		Total: 2,
		Results: []Result{
			{
				Module:  "github.com/spf13/cobra",
				Version: "v1.8.0",
			},
			{
				Module:  "github.com/broken/pkg",
				Version: "v1.0.0",
				Vulnerabilities: []Vulnerability{
					{ID: "GO-2024-1234", Summary: "RCE bug", Severity: "HIGH"},
				},
				Used: &used,
			},
		},
		Vulnerable: 1,
	}

	out, err := r.Render(s)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	// Проверяем, что в выводе есть ключевые подстроки.
	// Мы НЕ проверяем точный текст с цветами — это хрупко
	// (ANSI-коды зависят от терминала).
	// Проверяем "содержит ли" — это устойчиво.
	checks := []struct {
		name      string
		substring string
	}{
		{"общее количество", "2 шт"},
		{"первый модуль", "github.com/spf13/cobra"},
		{"версия", "v1.8.0"},
		{"id уязвимости", "GO-2024-1234"},
		{"серьёзность", "HIGH"},
		{"итог по уязвимым", "1 из 2"},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(out, c.substring) {
				t.Errorf("вывод не содержит %q.\nВывод:\n%s", c.substring, out)
			}
		})
	}
}

// TestTextRenderer_Empty — отдельный кейс: пустой ввод.
func TestTextRenderer_Empty(t *testing.T) {
	r := NewTextRenderer()
	out, err := r.Render(Summary{})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !strings.Contains(out, "Зависимости не найдены") {
		t.Errorf("ожидали сообщение о пустоте, got:\n%s", out)
	}
}

// TestHasVulnerabilities — проверка вспомогательного метода.
func TestHasVulnerabilities(t *testing.T) {
	tests := []struct {
		name string
		r    Result
		want bool
	}{
		{"пусто", Result{}, false},
		{"есть уязвимость", Result{Vulnerabilities: []Vulnerability{{}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.HasVulnerabilities(); got != tt.want {
				t.Errorf("HasVulnerabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}
