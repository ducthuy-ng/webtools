package webtools

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func ParseTemplates(templateDirectory string) (*template.Template, error) {
	dirFS := os.DirFS(templateDirectory)

	return ParseFS(dirFS, ".")
}

// ParseFS was written to support [embed.FS]
func ParseFS(templateFS fs.FS, root string) (*template.Template, error) {
	var parsedTemplates = template.Must(template.New("").Parse(""))

	err := fs.WalkDir(templateFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_ = fmt.Errorf("failed to parse path %s: %v", path, err)
		}

		if d.IsDir() {
			return nil
		}

		file, err := templateFS.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}

		fileStat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file size for reading: %v", err)
		}

		templateContent := make([]byte, fileStat.Size())
		_, err = file.Read(templateContent)
		if err != nil {
			return fmt.Errorf("failed to read file content: %v", err)
		}

		relPath, err := filepath.Rel(root, path)
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
