package views

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
)

func TestParseFSReturnsPageWithCategory(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"blog/index.tmpl": {
			Data: []byte("Category: {{category}}"),
		},
	}

	page, err := ParseFS(fsys, "blog/index.tmpl")
	if err != nil {
		t.Fatalf("ParseFS returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	page.Execute(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("Execute returned status %d, want %d", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if got := rr.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("Execute set Content-Type %q, want %q", got, "text/html; charset=utf-8")
	}

	if got := rr.Header().Get("Content-Length"); got != strconv.Itoa(len(body)) {
		t.Fatalf("Execute set Content-Length %q, want %q", got, strconv.Itoa(len(body)))
	}

	if !strings.Contains(body, "Category: blog") {
		t.Fatalf("Execute rendered %q, want to contain %q", body, "Category: blog")
	}
}

func TestParseFSReturnsErrorOnInvalidTemplate(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"blog/index.tmpl": {
			Data: []byte("{{define"),
		},
	}

	if _, err := ParseFS(fsys, "blog/index.tmpl"); err == nil {
		t.Fatalf("ParseFS expected to return an error for invalid template")
	}
}

func TestPageExecuteHandlesTemplateErrors(t *testing.T) {
	t.Parallel()

	tpl := template.Must(template.New("index.tmpl").Option("missingkey=error").Parse("Value: {{.Value}}"))
	page := Page{category: "blog", htmlTpl: tpl}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	page.Execute(rr, req, struct{}{})

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("Execute returned status %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "There was an error executing the template.") {
		t.Fatalf("Execute rendered %q, want error message", body)
	}
}

func TestPageExecuteHandlesHEADRequests(t *testing.T) {
	t.Parallel()

	const tmplBody = "Hello, world!"

	fsys := fstest.MapFS{
		"utils/index.tmpl": {
			Data: []byte(tmplBody),
		},
	}

	page, err := ParseFS(fsys, "utils/index.tmpl")
	if err != nil {
		t.Fatalf("ParseFS returned unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodHead, "/utils", nil)
	rr := httptest.NewRecorder()

	page.Execute(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("Execute returned status %d, want %d", rr.Code, http.StatusOK)
	}

	if rr.Body.Len() != 0 {
		t.Fatalf("Execute wrote %d bytes to HEAD response body, want 0", rr.Body.Len())
	}

	if got := rr.Header().Get("Content-Length"); got != strconv.Itoa(len(tmplBody)) {
		t.Fatalf("Execute set Content-Length %q, want %q", got, strconv.Itoa(len(tmplBody)))
	}
}
