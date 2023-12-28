package controllers

import (
	"github.com/go-chi/chi/v5"
)

var RestRouter = func(restApi chi.Router) {
	restApi.Route("/object", func(api chi.Router) {
		// FIXME: post-req body not passed??
		api.Post("/", GetSignedPutUrl)
		api.Get("/{file_key}", GetSignedGetUrl)
	})
	restApi.Route("/record", func(api chi.Router) {
		api.Post("/", CreateRecord)
		api.Get("/{file_key}", GetRecord)
	})
}
