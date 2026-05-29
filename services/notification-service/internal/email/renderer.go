package email

import (
	"bytes"
	"html/template"
	"path/filepath"
)

type Renderer struct {
	templateDir string
}

func NewRenderer(templateDir string) *Renderer {
	return &Renderer{
		templateDir: templateDir,
	}
}

func (r *Renderer) Render(templateName string, data any) (string, error) {
	files := []string{
		filepath.Join(
			r.templateDir,
			"layouts/base.html",
		),

		filepath.Join(
			r.templateDir,
			"partials/header.html",
		),

		filepath.Join(
			r.templateDir,
			"partials/footer.html",
		),

		filepath.Join(
			r.templateDir,
			templateName,
		),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	if err := tmpl.ExecuteTemplate(
		&buffer,
		"base",
		data,
	); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
