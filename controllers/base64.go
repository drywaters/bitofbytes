package controllers

import (
	"net/http"

	"github.com/DryWaters/bitofbytes/views"
)

type Base64Templates struct {
	Base64Response views.Page
}

func (u *Utils) Encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	encoded := u.Base64Service.Encode([]byte(r.PostFormValue("str")))
	e := struct {
		Response string
	}{
		Response: encoded,
	}
	u.Templates.Base64.Base64Response.Execute(w, r, e)
}

func (u *Utils) Decode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	str, err := u.Base64Service.Decode([]byte(r.PostFormValue("str")))
	if err != nil {
		http.Error(w, "unable to decode string", http.StatusBadRequest)
	}

	d := struct {
		Response string
	}{
		Response: str,
	}
	u.Templates.Base64.Base64Response.Execute(w, r, d)
}
