package controllers

import "net/http"

type Page interface {
	Execute(w http.ResponseWriter, r *http.Request, data any)
}
