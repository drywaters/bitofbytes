package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

func newTestBlogTemplates(t *testing.T) BlogTemplates {
	t.Helper()

	fsys := fstest.MapFS{
		"blog/index.tmpl": {
			Data: []byte("blog index"),
		},
		"blog/post.tmpl": {
			Data: []byte("Title: {{.Title}}\nContent: {{.Content}}"),
		},
	}

	index, err := views.ParseFS(fsys, "blog/index.tmpl")
	if err != nil {
		t.Fatalf("parse index template: %v", err)
	}
	post, err := views.ParseFS(fsys, "blog/post.tmpl")
	if err != nil {
		t.Fatalf("parse post template: %v", err)
	}

	return BlogTemplates{
		Index: index,
		Post:  post,
	}
}

func TestBlogIndexRendersTemplate(t *testing.T) {
	t.Parallel()

	blog := Blog{
		Templates: newTestBlogTemplates(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/blog", nil)
	rr := httptest.NewRecorder()

	blog.Index(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Index status code = %d, want %d", rr.Code, http.StatusOK)
	}

	if got := rr.Body.String(); got != "blog index" {
		t.Fatalf("expected template output %q, got %q", "blog index", got)
	}
}

func TestBlogReturnsNotFoundForMissingPost(t *testing.T) {
	t.Parallel()

	blog := Blog{
		PostService: models.PostService{},
		Templates:   newTestBlogTemplates(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/missing", nil)
	req.SetPathValue("slug", "missing")
	rr := httptest.NewRecorder()

	blog.Blog(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Blog status code = %d, want %d", rr.Code, http.StatusNotFound)
	}

	if got := strings.TrimSpace(rr.Body.String()); got != "Post not found" {
		t.Fatalf("expected not found message, got %q", got)
	}
}

func TestBlogRendersPostContent(t *testing.T) {
	t.Parallel()

	blog := Blog{
		PostService: models.PostService{},
		Templates:   newTestBlogTemplates(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/1", nil)
	req.SetPathValue("slug", "1")
	rr := httptest.NewRecorder()

	blog.Blog(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Blog status code = %d, want %d", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Title: Getting Started") {
		t.Fatalf("expected post title in body, got %q", body)
	}
	if !strings.Contains(body, "Warming up") {
		t.Fatalf("expected post content in body, got %q", body)
	}
}
