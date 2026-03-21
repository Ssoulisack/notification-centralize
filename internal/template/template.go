package template

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/your-org/notification-center/internal/model"
)

// Engine renders notification templates with data.
type Engine struct{}

func New() *Engine {
	return &Engine{}
}

// RenderSubject renders the subject line using text/template (no HTML).
func (e *Engine) RenderSubject(tmpl *model.Template, data map[string]string) (string, error) {
	t, err := texttemplate.New("subject").Parse(tmpl.SubjectTemplate)
	if err != nil {
		return "", fmt.Errorf("parse subject template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute subject template: %w", err)
	}
	return buf.String(), nil
}

// RenderBody renders the body using html/template (safe HTML output).
func (e *Engine) RenderBody(tmpl *model.Template, data map[string]string) (string, error) {
	t, err := htmltemplate.New("body").Parse(tmpl.BodyTemplate)
	if err != nil {
		return "", fmt.Errorf("parse body template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute body template: %w", err)
	}
	return buf.String(), nil
}
