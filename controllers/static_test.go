package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/DryWaters/bitofbytes/views"
)

func TestStaticHandlerExecutesTemplate(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		"static/page.tmpl": {
			Data: []byte("static page"),
		},
	}

	page, err := views.ParseFS(fsys, "static/page.tmpl")
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}

	handler := StaticHandler(page)

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("StaticHandler status code = %d, want %d", rr.Code, http.StatusOK)
	}

	if got := rr.Body.String(); got != "static page" {
		t.Fatalf("expected template output %q, got %q", "static page", got)
	}
}
