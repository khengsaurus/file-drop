package controllers

import (
	"github.com/go-chi/chi/v5"
)

var RestRouter = func(restApi chi.Router) {
	restApi.Route("/object", func(api chi.Router) {
		api.Post("/", GetSignedPutUrl)
		api.Get("/{file_key}", GetSignedGetUrl)
	})
	restApi.Route("/record", func(api chi.Router) {
		api.Post("/", CreateRecord)
		api.Get("/{file_key}", GetRecord)
	})
}

var AdminRouter = func(adminRouter chi.Router) {
	adminRouter.Get("/", AdminGet)
	adminRouter.Delete("/{action}", AdminDelete)
}
