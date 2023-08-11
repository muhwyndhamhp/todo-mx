package public

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplateHandler(e *echo.Echo) {
	tmpl := template.Must(template.ParseGlob("public/views/*.html"))
	template.Must(tmpl.ParseGlob("public/components/*.html"))
	template.Must(tmpl.ParseGlob("public/templates/*.html"))

	t := newTemplate(tmpl)
	e.Renderer = t
}

func newTemplate(templates *template.Template) echo.Renderer {
	return &Template{
		Templates: templates,
	}
}
