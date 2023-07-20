package templates

import (
	"fmt"
	"html/template"
)

type MainPage struct {
	OrderExists bool
	OrderJSON   []byte
	Message     string
}

func GetMainTemplate() (*template.Template, error) {
	const op = "web.templates.GetMainTemplate"

	tmpl, err := template.ParseFiles("/Users/nikolay/Projects/WB/internship/L0/Project/internal/web/templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return tmpl, nil
}
