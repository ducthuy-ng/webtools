package webtools

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func ParseTemplates(templateDirectory string) (*template.Template, error) {
	var parsedTemplates = template.Must(template.New("/").Parse(""))

	err := filepath.WalkDir(templateDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_ = fmt.Errorf("failed to parse path %s: %v", path, err)
		}

		if d.IsDir() {
			return nil
		}

		templateContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file content: %v", err)
		}

		relPath, err := filepath.Rel(templateDirectory, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		template.Must(parsedTemplates.New(relPath).Parse(string(templateContent)))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk dir: %v", err)
	}
	return parsedTemplates, nil
}
