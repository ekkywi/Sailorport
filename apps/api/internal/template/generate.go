package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GenerateInput struct {
	TemplateDir string
	TargetDir   string
	Name        string
}

func Generate(in GenerateInput) error {
	if err := os.MkdirAll(in.TargetDir, 0o755); err != nil {
		return fmt.Errorf("Create target directory: %w", err)
	}

	return filepath.WalkDir(in.TemplateDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(in.TemplateDir, path)
		if err != nil {
			return err
		}
		if rel == "." || rel == "manifest.json" {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(in.TargetDir, rel), 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		outName := rel
		content := string(data)
		if strings.HasSuffix(rel, ".tmpl") {
			outName = strings.TrimSuffix(rel, ".tmpl")
			tmpl, err := template.New(outName).Parse(content)
			if err != nil {
				return fmt.Errorf("Parse template %s: %w", rel, err)
			}
			var buf strings.Builder
			if err := tmpl.Execute(&buf, map[string]string{"Name": in.Name}); err != nil {
				return err
			}
			content = buf.String()
		}

		outPath := filepath.Join(in.TargetDir, outName)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(outPath, []byte(content), 0o644)
	})
}
