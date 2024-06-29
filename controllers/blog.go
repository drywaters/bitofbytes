package controllers

import (
	"bytes"
	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
	"github.com/adrg/frontmatter"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"net/http"
	"strings"
)

type Blog struct {
	PostService models.PostService
	Templates   BlogTemplates
}

type BlogTemplates struct {
	Index views.Page
	Post  views.Page
}

func (b Blog) Index(w http.ResponseWriter, r *http.Request) {
	b.Templates.Index.Execute(w, r, nil)
}

func (b Blog) Blog(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	slug := chi.URLParam(r, "slug")
	postMarkdown, err := b.PostService.Read(slug)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	rest, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)
	var buf bytes.Buffer
	err = mdRenderer.Convert(rest, &buf)
	if err != nil {
		http.Error(w, "Error converting markdown", http.StatusInternalServerError)
		return
	}
	post.Content = template.HTML(buf.String())
	b.Templates.Post.Execute(w, r, post)
}
