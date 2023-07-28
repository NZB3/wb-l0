package templates

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type MainPage struct {
	OrderExists bool
	OrderJSON   []byte
	Message     string
}

func GetMainTemplate() (*template.Template, error) {
	const op = "web.templates.GetMainTemplate"
	htmlPath, err := filepath.Abs("../../internal/web/templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return tmpl, nil
}
