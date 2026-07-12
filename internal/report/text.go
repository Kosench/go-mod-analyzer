package report

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")). // светлый текст
	Background(lipgloss.Color("#7D56F4")). // фиолетовый фон
	Padding(0, 1)

var normalStyle = lipgloss.NewStyle()

// Indirect-зависимости рисуем тусклым (серым) цветом — визуально отделяем.
var dimStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#626262")) // тёмно-серый

// Метка уязвимости: красный, жирный — привлекает внимание.
var vulnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000")).
	Bold(true)

// Renderer — интерфейс форматтера вывода.
//
// Это позволит нам иметь разные реализации (text, json, sarif, html)
// с единым интерфейсом. main.go выберет нужный по флагу --format.

type Renderer interface {
	Render(s Summary) (string, error)
}

type TextRenderer struct{}

func NewTextRenderer() *TextRenderer {
	return &TextRenderer{}
}

func (r *TextRenderer) Render(s Summary) (string, error) {
	var b strings.Builder

	// Заголовок с общей статистикой.
	b.WriteString(titleStyle.Render(fmt.Sprintf(
		" Анализ зависимостей: %d шт. ", s.Total,
	)))
	b.WriteString("\n\n")

	// Если результатов нет — ранний возврат.
	// Избегаем глубокой вложенности (early return — хорошая практика).
	if len(s.Results) == 0 {
		b.WriteString(dimStyle.Render("Зависимости не найдены.\n"))
		return b.String(), nil
	}

	// Проходим по каждому результату и формируем строку.
	for _, res := range s.Results {
		line := fmt.Sprintf("  • %s %s", res.Module, res.Version)

		// Indirect помечаем серым.
		if res.Indirect {
			line = dimStyle.Render(line + " (indirect)")
		} else {
			line = normalStyle.Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")

		// Уязвимости — отдельными строками под модулем.
		for _, v := range res.Vulnerabilities {
			vulnLine := fmt.Sprintf("      ⚠ %s [%s] %s", v.ID, v.Severity, v.Summary)
			b.WriteString(vulnStyle.Render(vulnLine))
			b.WriteString("\n")
		}
	}

	// Итоговая строка про уязвимости.
	if s.Vulnerable > 0 {
		b.WriteString("\n")
		b.WriteString(vulnStyle.Render(fmt.Sprintf(
			"Найдено уязвимостей в %d из %d зависимостей.",
			s.Vulnerable, s.Total,
		)))
		b.WriteString("\n")
	}

	return b.String(), nil
}
