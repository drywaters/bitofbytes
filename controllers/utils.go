package controllers

import (
	"net/http"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

type Utils struct {
	Base64Service models.Base64Service
	Templates     UtilsTemplates
}

type UtilsTemplates struct {
	Index  views.Page
	Base64 Base64Templates
}

func (u *Utils) Index(w http.ResponseWriter, r *http.Request) {
	u.Templates.Index.Execute(w, r, nil)
}
