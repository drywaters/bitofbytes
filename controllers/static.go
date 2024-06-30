package controllers

import (
	"github.com/DryWaters/bitofbytes/views"
	"net/http"
)

func StaticHandler(p views.Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.Execute(w, r, nil)
	}
}
