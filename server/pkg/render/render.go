package render

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var templatePath, _ = filepath.Abs("./")
	templatePath = filepath.Join(templatePath, "server", "templates")
	var err error
	t.Templates, err = template.ParseFiles(filepath.Join(templatePath, "base.layout.tmpl"),
		filepath.Join(templatePath, name))
	if err != nil {
		return err
	}
	return t.Templates.ExecuteTemplate(w, name, data)
}
