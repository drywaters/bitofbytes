package controllers

import (
	"net/http"

	"github.com/DryWaters/bitofbytes/views"
)

type StringsTemplates struct {
	LengthResponse views.Page
}

func (u *Utils) StringsLength(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	details := u.UtilsService.StringDetails(r.PostFormValue("str"))

	u.Templates.Strings.LengthResponse.Execute(w, r, details)
}
