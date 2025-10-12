package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

const csrfCookiePath = "/"

// CSRF returns middleware that configures gorilla/csrf with the provided key and flags.
func CSRF(key string, secure bool) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(key),
		csrf.Secure(secure),
		csrf.Path(csrfCookiePath),
	)
}
