package views

import (
	"io"
	"text/template"
)

type Templates struct {
	T *template.Template
}

// Render template
func (t *Templates) Execute(rw io.Writer, name string, data any) error {
	return t.T.ExecuteTemplate(rw, name, data)
}
