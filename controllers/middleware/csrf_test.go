package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"
)

func TestCSRF(t *testing.T) {
	t.Parallel()

	const key = "01234567890123456789012345678901"

	handler := CSRF(key, true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := csrf.Token(r); token == "" {
			t.Fatal("expected CSRF token to be present on the request context")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	var csrfCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == "_gorilla_csrf" {
			csrfCookie = c
			break
		}
	}

	if csrfCookie == nil {
		t.Fatalf("expected CSRF cookie to be set")
	}

	if !csrfCookie.Secure {
		t.Errorf("expected CSRF cookie to be secure")
	}

	if csrfCookie.Path != "/" {
		t.Errorf("expected CSRF cookie path to be '/', got %q", csrfCookie.Path)
	}
}
