package controllers

import (
	"net/http"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

type Utils struct {
	UtilsService models.UtilsService
	Templates    UtilsTemplates
}

type UtilsTemplates struct {
	Index   views.Page
	Strings StringsTemplates
}

func (u *Utils) Index(w http.ResponseWriter, r *http.Request) {
	u.Templates.Index.Execute(w, r, nil)
}
