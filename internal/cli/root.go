// Package cli handles command-line argument parsing.
//
// It contains the root Cobra command. It does NOT contain any business logic —
// it merely collects configuration from flags and passes it to analyzer/report.
package cli

import (
	"fmt"
	"io"

	"github.com/Kosench/go-mod-analyzer/internal/module"
	"github.com/Kosench/go-mod-analyzer/internal/report"
	"github.com/spf13/cobra"
)

// Options — конфигурация запуска, собранная из флагов и аргументов.
// Options — конфигурация запуска, собранная из флагов и аргументов.
//
// Передаётся как одна структура, а не россыпью аргументов:
// это удобнее расширять (добавил поле — добавил флаг).
type Options struct {
	// Path — путь к проекту (аргумент, не флаг).
	Path string

	// Format — формат вывода: text, json, sarif, html, tui.
	Format string

	// IndirectOnly — показывать только indirect (пример будущего флага).
	// Пока не используется, но показывает, как расширять Options.
}

// rootCmd — главная команда.
// Создаём её в переменной, чтобы cobra мог её найти.
//
// Делаем НЕэкспортируемое поле и фабрику NewRootCmd — это позволяет
// внедрять зависимости (writer для вывода) и тестировать.
func NewRootCmd(out io.Writer) *cobra.Command {
	// Конфигурация, которая накопится из флагов.
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "go-mod-analyzer [path]",
		Short: "Анализ зависимостей Go-проекта",
		Long: "go-mod-analyzer — утилита для анализа зависимостей: поиск уязвимостей\n" +
			"(через OSV API), выявление неиспользуемых модулей и построение графа.",
		// Args указывает cobra, какие аргументы ожидать.
		// MaximumNArgs(1): не больше одного позиционного аргумента (путь).
		Args: cobra.MaximumNArgs(1),
		// RunE — функция, которая выполняется при запуске команды.
		// Возвращает error (в отличие от Run), что позволяет cobra
		// корректно показать ошибку и выставить exit code.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Если путь передан аргументом — берём его.
			// Иначе по умолчанию текущая директория.
			if len(args) > 0 {
				opts.Path = args[0]
			} else {
				opts.Path = "."
			}

			return runAnalyze(cmd, out, opts)
		},
	}

	// Регистрируем флаги.
	// cmd.Flags() — локальные флаги этой команды.
	cmd.Flags().StringVarP(
		&opts.Format, // куда записать значение
		"format",     // длинное имя
		"f",          // короткое имя (-f)
		"text",       // значение по умолчанию
		"формат вывода: text|json|sarif|html|tui",
	)

	return cmd
}

// runAnalyze — оркестрация: парсинг → построение отчёта → рендер.
//
// Вынесена отдельно, чтобы её можно было покрыть тестами.
// Принимает io.Writer — это позволяет писать в bytes.Buffer в тестах.
func runAnalyze(cmd *cobra.Command, out io.Writer, opts Options) error {
	// 1. Парсим go.mod.
	modules, err := module.ParseFile(opts.Path)
	if err != nil {
		// Возвращаем ошибку наверх — cobra сам напечатает её в stderr.
		return fmt.Errorf("парсинг go.mod: %w", err)
	}

	// 2. Строим Summary из модулей.
	// Пока без уязвимостей и without unused-анализа — это следующие этапы.
	summary := buildSummary(modules)

	// 3. Выбираем рендерер по формату.
	renderer, err := selectRenderer(opts.Format)
	if err != nil {
		return err
	}

	// 4. Рендерим и пишем в out.
	output, err := renderer.Render(summary)
	if err != nil {
		return fmt.Errorf("рендер: %w", err)
	}

	// Fprintln добавляет перенос в конце — терминальные утилиты так делают.
	fmt.Fprintln(out, output)

	return nil
}

// buildSummary преобразует []Module в report.Summary.
// Здесь — место, куда на следующих этапах добавится проверка уязвимостей
// и поиск unused-пакетов.
func buildSummary(modules []module.Module) report.Summary {
	results := make([]report.Result, 0, len(modules))
	vulnerable := 0

	for _, m := range modules {
		results = append(results, report.Result{
			Module:   m.Path,
			Version:  m.Version,
			Indirect: m.Indirect,
		})
	}

	return report.Summary{
		Results:    results,
		Total:      len(modules),
		Vulnerable: vulnerable, // пока всегда 0
	}
}

// selectRenderer выбирает форматтер по имени формата.
//
// Это "фабрика" — простой switch. Когда форматов станет много,
// можно будет вынести в map[string]Renderer или реестр.
func selectRenderer(format string) (report.Renderer, error) {
	switch format {
	case "text":
		return report.NewTextRenderer(), nil
	default:
		// fmt.Errorf с %q обрамляет значение кавычками — нагляднее в ошибке.
		return nil, fmt.Errorf("неподдерживаемый формат %q (доступен: text)", format)
	}
}
