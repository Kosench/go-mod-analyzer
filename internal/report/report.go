// Package report is responsible for presenting analysis results.
//
// It contains result types and output formatters (text, json, sarif, html, tui).
// The package does NOT know where the data came from — it only displays them.
package report

type Result struct {
	Module  string
	Version string
	// Indirect true для транзитивных зависимостей.
	Indirect        bool
	Vulnerabilities []Vulnerability
	// Used true, если модуль реально импортируется в коде.
	// Будет заполняться на Этапе 3 (AST-анализ).
	// Пока nil-семантика: статус "неизвестен".
	Used *bool
}

type Vulnerability struct {
	ID       string
	Summary  string
	Severity string
}

func (r Result) HasVulnerabilities() bool {
	return len(r.Vulnerabilities) > 0
}

type Summary struct {
	Results    []Result
	Total      int
	Vulnerable int
}
