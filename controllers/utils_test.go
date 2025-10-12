package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/DryWaters/bitofbytes/views"
)

func newTestIndexPage(t *testing.T, body string) views.Page {
	t.Helper()

	fsys := fstest.MapFS{
		"utils/index.tmpl": {
			Data: []byte(body),
		},
	}

	page, err := views.ParseFS(fsys, "utils/index.tmpl")
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}

	return page
}

func TestUtilsIndexRendersTemplate(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Templates: UtilsTemplates{
			Index: newTestIndexPage(t, "utility index"),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/utils", nil)
	rr := httptest.NewRecorder()

	utils.Index(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Index status code = %d, want %d", rr.Code, http.StatusOK)
	}

	if got := rr.Body.String(); got != "utility index" {
		t.Fatalf("expected template output %q, got %q", "utility index", got)
	}
}
