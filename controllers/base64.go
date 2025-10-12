package controllers

import (
	"errors"
	"net/http"

	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/views"
)

type Base64Templates struct {
	Base64Response views.Page
}

type base64Response struct {
	Response string
	Error    string
}

func (u *Utils) Encode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	encoded := u.Base64Service.Encode([]byte(r.PostFormValue("str")))
	u.Templates.Base64.Base64Response.Execute(w, r, base64Response{
		Response: encoded,
	})
}

func (u *Utils) Decode(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	str, err := u.Base64Service.Decode([]byte(r.PostFormValue("str")))
	if err != nil {
		message := "We couldn't decode the submitted text. Please verify it is valid Base64."
		if errors.Is(err, models.ErrBase64InputTooLarge) {
			message = "The submitted text is too large to decode. Please limit it to 4 KB or less."
		}

		u.Templates.Base64.Base64Response.Execute(w, r, base64Response{
			Error: message,
		})
		return
	}

	u.Templates.Base64.Base64Response.Execute(w, r, base64Response{
		Response: str,
	})
}
