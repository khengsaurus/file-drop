package controllers

import (
	"github.com/go-chi/chi/v5"
)

type Payload struct {
	FileName string `json:"fileName"`
	Key      string `json:"key"`
	Url      string `json:"url"`
}

var RestRouter = func(restApi chi.Router) {
	restApi.Route("/file", func(fileApi chi.Router) {
		fileApi.Post("/", GetSignedPutUrl)
		fileApi.Get("/{file_key}", GetSignedGetUrl)
		// fileApi.Delete("/{file_key}", DeleteFile)
	})
	restApi.Route("/record", func(recordApi chi.Router) {
		recordApi.Post("/", CreateRecord)
		recordApi.Get("/{file_key}", GetRecord)
	})
}

var AdminRouter = func(adminRouter chi.Router) {
	adminRouter.Get("/", AdminGet)
	adminRouter.Delete("/{action}", AdminDelete)
}
