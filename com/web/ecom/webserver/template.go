package webserver

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo"

	"o.o/backend/pkg/common/cmenv"
)

type Templates struct {
	templates   *template.Template
	filePattern string
	autoReload  bool
	lastReload  time.Time
}

func parseTemplates(pattern string) (*Templates, error) {
	templates, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	t := &Templates{
		templates:   templates,
		filePattern: pattern,
		autoReload:  cmenv.IsDev(),
	}
	return t, nil
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.autoReload && time.Now().Sub(t.lastReload) > 5*time.Second {
		templates, err := template.ParseGlob(t.filePattern)
		if err != nil {
			return err
		}
		t.lastReload, t.templates = time.Now(), templates
		return templates.ExecuteTemplate(w, name, data)
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *Templates) Reload() error {
	templates, err := template.ParseGlob(t.filePattern)
	if err != nil {
		return err
	}
	t.templates = templates
	return nil
}
