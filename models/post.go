package models

import (
	"html/template"
	"io/fs"

	"github.com/DryWaters/bitofbytes/posts"
)

type Post struct {
	Title   string `toml:"title"`
	Slug    string `toml:"slug"`
	Content template.HTML
	Author  Author `toml:"author"`
}

type Author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type PostService struct {
}

func (p PostService) Read(slug string) (string, error) {
	b, err := fs.ReadFile(posts.FS, slug+".md")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
