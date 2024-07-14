package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
)

func Must(p Page, err error) Page {
	if err != nil {
		panic(err)
	}
	return p
}

func ParseFS(fs fs.FS, patterns ...string) (Page, error) {
	var page Page
	tpl, err := template.New(path.Base(patterns[0])).Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
		},
	).ParseFS(fs, patterns...)

	if err != nil {
		return page, fmt.Errorf("parsing template: %w", err)
	}
	return Page{
		htmlTpl: tpl,
		name:    path.Dir(patterns[0]),
	}, nil
}

type Page struct {
	name    string
	htmlTpl *template.Template
}

type PageData map[string]any

func (p Page) Execute(w http.ResponseWriter, r *http.Request, data any) {
	pageData := struct {
		PageData PageData
		Data     any
	}{
		PageData: map[string]any{"Name": p.name},
		Data:     data,
	}
	tpl, err := p.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})
	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, pageData)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}
