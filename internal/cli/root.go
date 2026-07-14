// Package cli handles command-line argument parsing.
//
// It contains the root Cobra command. It does NOT contain any business logic —
// it merely collects configuration from flags and passes it to analyzer/report.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Kosench/go-mod-analyzer/internal/module"
	"github.com/Kosench/go-mod-analyzer/internal/report"
	"github.com/spf13/cobra"
)

// Options — конфигурация запуска.
type Options struct {
	Path   string
	Format string
}

func NewRootCmd(out io.Writer, parser module.Parser, version string) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "go-mod-analyzer [path]",
		Short: "Analyze Go project dependencies",
		Long:  "go-mod-analyzer — utility for analyzing dependencies...",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = args[0]
			} else {
				opts.Path = "."
			}
			return runAnalyze(cmd, out, opts, parser, version)
		},
	}

	cmd.Flags().StringVarP(&opts.Format, "format", "f", "text", "output format: text|json|sarif|html|tui")
	return cmd
}

func resolveModPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat path %q: %w", path, err)
	}
	if info.IsDir() {
		return filepath.Join(path, "go.mod"), nil
	}
	return path, nil
}

func runAnalyze(cmd *cobra.Command, out io.Writer, opts *Options, parser module.Parser, version string) error {
	modPath, err := resolveModPath(opts.Path)
	if err != nil {
		return err
	}

	modules, err := parser.Parse(modPath)
	if err != nil {
		// ИСПРАВЛЕНО: go.mod
		return fmt.Errorf("parse go.mod: %w", err)
	}

	summary := buildSummary(modules)

	renderer, err := selectRenderer(opts.Format, version)
	if err != nil {
		return err
	}

	output, err := renderer.Render(summary)
	if err != nil {
		return fmt.Errorf("render report: %w", err)
	}

	fmt.Fprintln(out, output)
	return nil
}

func buildSummary(modules []module.Module) report.Summary {
	results := make([]report.Result, 0, len(modules))
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
		Vulnerable: 0,
	}
}

func selectRenderer(format, version string) (report.Renderer, error) {
	switch format {
	case "text":
		return report.NewTextRenderer(), nil
	default:
		return nil, fmt.Errorf("unsupported format %q (available: text)", format)
	}
}
