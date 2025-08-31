package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

type Utils struct {
	Base64Service models.Base64Service
	Templates     UtilsTemplates
}

type UtilsTemplates struct {
	Index            views.Page
	Base64           Base64Templates
	Base64EncodePage views.Page
	Base64DecodePage views.Page
}

func (u *Utils) Index(w http.ResponseWriter, r *http.Request) {
	u.Templates.Index.Execute(w, r, nil)
}

// EncodeGET serves the encode page by default, or returns an encoded result
// when the query parameter `q` is provided.
func (u *Utils) EncodeGET(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		u.Templates.Base64EncodePage.Execute(w, r, nil)
		return
	}

	if len(q) > maxBase64InputLen {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "input too long", "message": "Input too long (max 1024 chars)"})
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("Input too long (max 1024 chars)"))
		return
	}

	urlSafe := parseBool(r.URL.Query().Get("urlSafe"))
	noPad := parseBool(r.URL.Query().Get("noPad"))
	res := u.Base64Service.EncodeWith([]byte(q), urlSafe, noPad)
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": res})
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(res))
}

// DecodeGET serves the decode page by default, or returns a decoded result
// when the query parameter `q` is provided.
func (u *Utils) DecodeGET(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		u.Templates.Base64DecodePage.Execute(w, r, nil)
		return
	}

	if len(q) > maxBase64InputLen {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(map[string]string{"error": "input too long", "message": "Input too long (max 1024 chars)"})
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("Input too long (max 1024 chars)"))
		return
	}

	urlSafe := parseBool(r.URL.Query().Get("urlSafe"))
	noPad := parseBool(r.URL.Query().Get("noPad"))
	str, err := u.Base64Service.DecodeWith([]byte(q), urlSafe, noPad)
	if err != nil {
		if strings.Contains(r.Header.Get("Accept"), "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "unable to decode", "message": "Unable to decode string"})
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to decode string"))
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": str})
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(str))
}

func parseBool(v string) bool {
	switch strings.ToLower(v) {
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}
