package controllers

import "github.com/DryWaters/bitofbytes/views"
import "github.com/DryWaters/bitofbytes/models"

type Utils struct {
	UtilsService models.UtilsService
	Templates    UtilsTemplates
}

type UtilsTemplates struct {
	Index views.Page
}
