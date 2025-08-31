package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DryWaters/bitofbytes/views"
)

const maxBase64InputLen = 1024

type Base64Templates struct {
	Base64EncodeResponse views.Page
	Base64DecodeResponse views.Page
	Base64Error          views.Page
}

func (u *Utils) Encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	in := r.PostFormValue("str")
	urlSafe := r.PostFormValue("urlSafe") != ""
	noPad := r.PostFormValue("noPad") != ""
	if len(in) > maxBase64InputLen {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "input too long", "message": "Input too long (max 1024 chars)"})
		} else {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			u.Templates.Base64.Base64Error.Execute(w, r, struct{ Error string }{Error: "Input too long (max 1024 chars)"})
		}
		return
	}
	encoded := u.Base64Service.EncodeWith([]byte(in), urlSafe, noPad)
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": encoded})
		return
	}
	e := struct{ Response string }{Response: encoded}
	u.Templates.Base64.Base64EncodeResponse.Execute(w, r, e)
}

func (u *Utils) Decode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	in := r.PostFormValue("str")
	urlSafe := r.PostFormValue("urlSafe") != ""
	noPad := r.PostFormValue("noPad") != ""
	if len(in) > maxBase64InputLen {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "input too long", "message": "Input too long (max 1024 chars)"})
		} else {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			u.Templates.Base64.Base64Error.Execute(w, r, struct{ Error string }{Error: "Input too long (max 1024 chars)"})
		}
		return
	}

	str, err := u.Base64Service.DecodeWith([]byte(in), urlSafe, noPad)
	if err != nil {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "unable to decode", "message": "Unable to decode string"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			u.Templates.Base64.Base64Error.Execute(w, r, struct{ Error string }{Error: "Unable to decode string"})
		}
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": str})
		return
	}
	d := struct{ Response string }{Response: str}
	u.Templates.Base64.Base64DecodeResponse.Execute(w, r, d)
}
