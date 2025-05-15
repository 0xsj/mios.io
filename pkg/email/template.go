package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type TemplateManager struct {
	templates map[string]*template.Template
	baseDir   string
}

func NewTemplateManager(baseDir string) (*TemplateManager, error) {
	manager := &TemplateManager{
		templates: make(map[string]*template.Template),
		baseDir:   baseDir,
	}

	// Load all templates
	err := manager.loadTemplates()
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *TemplateManager) loadTemplates() error {
	// Check if directory exists
	_, err := os.Stat(m.baseDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("template directory %s does not exist", m.baseDir)
	}

	// Walk through the directory
	return filepath.Walk(m.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .html or .txt files
		ext := filepath.Ext(path)
		if ext != ".html" && ext != ".txt" {
			return nil
		}

		// Read template file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse template
		name := filepath.Base(path)
		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return err
		}

		// Store template
		m.templates[name] = tmpl
		return nil
	})
}

func (m *TemplateManager) Render(name string, data interface{}) (string, error) {
	tmpl, exists := m.templates[name]
	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Add adds a template
func (m *TemplateManager) Add(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}

	m.templates[name] = tmpl
	return nil
}
