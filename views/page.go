package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
)

type Page struct {
	category string
	htmlTpl  *template.Template
}

func (p Page) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tpl, err := p.htmlTpl.Clone()
	if err != nil {
		slog.Error("cloning template", "error", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"category": func() string {
			return p.category
		},
	})
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		slog.Error("executing template", "error", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}

	content := buf.Bytes()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if _, err := w.Write(content); err != nil {
		slog.Error("writing template response", "error", err)
	}
}

func Must(p Page, err error) Page {
	if err != nil {
		panic(err)
	}
	return p
}

func ParseFS(fs fs.FS, patterns ...string) (Page, error) {
	var page Page
	category, _, _ := strings.Cut(patterns[0], "/")
	tpl, err := template.New(path.Base(patterns[0])).Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"category": func() string {
				return ""
			},
		},
	).ParseFS(fs, patterns...)

	if err != nil {
		return page, fmt.Errorf("parsing template: %w", err)
	}

	return Page{
		htmlTpl:  tpl,
		category: category,
	}, nil
}
