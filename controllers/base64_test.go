package controllers

import (
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

func newTestBase64Page(t *testing.T) views.Page {
	t.Helper()

	fsys := fstest.MapFS{
		"utils/base64_response.tmpl": {
			Data: []byte("Response: {{.Response}}\nError: {{.Error}}"),
		},
	}

	page, err := views.ParseFS(fsys, "utils/base64_response.tmpl")
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}

	return page
}

func TestUtilsEncodeWritesEncodedString(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Base64Service: models.Base64Service{},
		Templates: UtilsTemplates{
			Base64: Base64Templates{Base64Response: newTestBase64Page(t)},
		},
	}

	form := url.Values{}
	form.Set("encoding", "standard")
	form.Set("str", "hello world")
	req := httptest.NewRequest(http.MethodPost, "/encode", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	utils.Encode(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Encode status code = %d, want %d", rr.Code, http.StatusOK)
	}

	if got := rr.Header().Get("Content-Type"); got != "text/html" {
		t.Fatalf("Content-Type header = %q, want %q", got, "text/html")
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Response: aGVsbG8gd29ybGQ=") {
		t.Fatalf("expected encoded string in body, got %q", body)
	}
	if !strings.Contains(body, "Error: ") {
		t.Fatalf("expected empty error field in body, got %q", body)
	}
}

func TestUtilsDecodeWritesDecodedString(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Base64Service: models.Base64Service{},
		Templates: UtilsTemplates{
			Base64: Base64Templates{Base64Response: newTestBase64Page(t)},
		},
	}

	form := url.Values{}
	form.Set("encoding", "standard")
	form.Set("str", "aGVsbG8gd29ybGQ=")
	req := httptest.NewRequest(http.MethodPost, "/decode", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	utils.Decode(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Decode status code = %d, want %d", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Response: hello world") {
		t.Fatalf("expected decoded string in body, got %q", body)
	}
	if !strings.Contains(body, "Error: ") {
		t.Fatalf("expected empty error field in body, got %q", body)
	}
}

func TestUtilsDecodeHandlesInvalidInput(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Base64Service: models.Base64Service{},
		Templates: UtilsTemplates{
			Base64: Base64Templates{Base64Response: newTestBase64Page(t)},
		},
	}

	form := url.Values{}
	form.Set("encoding", "standard")
	form.Set("str", "not-base64")
	req := httptest.NewRequest(http.MethodPost, "/decode", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	utils.Decode(rr, req)

	body := html.UnescapeString(rr.Body.String())
	want := "We couldn't decode the submitted text. Please verify it is valid Base64."
	if !strings.Contains(body, want) {
		t.Fatalf("expected error message %q, got %q", want, body)
	}
}

func TestUtilsDecodeHandlesOversizedInput(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Base64Service: models.Base64Service{},
		Templates: UtilsTemplates{
			Base64: Base64Templates{Base64Response: newTestBase64Page(t)},
		},
	}

	largeInput := strings.Repeat("A", 4097)
	form := url.Values{}
	form.Set("encoding", "standard")
	form.Set("str", largeInput)
	req := httptest.NewRequest(http.MethodPost, "/decode", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	utils.Decode(rr, req)

	body := html.UnescapeString(rr.Body.String())
	want := "The submitted text is too large to decode. Please limit it to 4 KB or less."
	if !strings.Contains(body, want) {
		t.Fatalf("expected oversized error message %q, got %q", want, body)
	}
}

func TestUtilsHandlesUnknownEncoding(t *testing.T) {
	t.Parallel()

	utils := &Utils{
		Base64Service: models.Base64Service{},
		Templates: UtilsTemplates{
			Base64: Base64Templates{Base64Response: newTestBase64Page(t)},
		},
	}

	form := url.Values{}
	form.Set("encoding", "unknown")
	form.Set("str", "hello")

	req := httptest.NewRequest(http.MethodPost, "/encode", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	utils.Encode(rr, req)

	body := html.UnescapeString(rr.Body.String())
	want := "The selected Base64 variant is not supported."
	if !strings.Contains(body, want) {
		t.Fatalf("expected unknown encoding error message %q, got %q", want, body)
	}
}
