package controllers

import (
	"net/http"
)

func StaticHandler(p Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.Execute(w, r, nil)
	}
}
