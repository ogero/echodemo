package components

import (
	"bitbucket.org/ogero/echodemo/embed"
	"errors"
	"fmt"
	"github.com/Unknwon/i18n"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
	templates           map[string]*template.Template
}

func NewTemplateRenderer() (renderer *TemplateRenderer, err error) {
	renderer = &TemplateRenderer{
		TemplateLayoutPath:  "templates/layout/",
		TemplateIncludePath: "templates/views/",
	}
	if err = renderer.Load(); err != nil {
		renderer = nil
	}
	return
}

func (t *TemplateRenderer) Load() error {
	if t.templates == nil {
		t.templates = make(map[string]*template.Template)
	}
	var layoutFiles, includeFiles []string
	var err error

	if layoutFiles, err = embed.WalkDirs(t.TemplateLayoutPath, false); err != nil {
		return err
	}

	if includeFiles, err = embed.WalkDirs(t.TemplateIncludePath, false); err != nil {
		return err
	}

	mTemplate := template.New("")
	mTemplate, err = mTemplate.Parse(`{{define "default" }} {{ template "default-layout" . }} {{ end }}`)
	if err != nil {
		return err
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		t.templates[fileName], err = mTemplate.Clone()
		if err != nil {
			return err
		}
		t.templates[fileName] = template.Must(parseEmbedFiles(t.templates[fileName], files...))
	}
	return nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	cc := c.(*EchoContext)
	layout := "default"
	names := strings.Split(name, "/")
	if len(names) == 2 {
		layout = names[0]
		name = names[1]
	}
	tpl, ok := t.templates[name]
	if !ok {
		return errors.New(fmt.Sprintf("The template %s does not exist.", name))
	}

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["isTruePtr"] = func(v *bool) bool {
			return v != nil && (*v) == true
		}
		viewContext["isLoggedUser"] = cc.GetLoggedUserID() > 0
		viewContext["loggedUser"] = cc.GetLoggedUser()
		viewContext["reverse"] = cc.Echo().Reverse
		viewContext["menu"] = func() []Menu { return GetMenu(cc) }
		viewContext["Tr"] = func(format string, args ...interface{}) template.HTML {
			return template.HTML(i18n.Tr(cc.GetLang(), format, args))
		}
		viewContext["flashes"] = cc.Flashes
	}

	err := tpl.ExecuteTemplate(w, layout, data)
	if err != nil {
		logrus.Debug(err)
	}
	return err
}

func parseEmbedFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := embed.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
