package controllers

import (
	"github.com/go-chi/chi/v5"
)

var ApiRouter = func(router chi.Router) {
	router.Route("/object", func(api chi.Router) {
		api.Post("/", GetSignedPutUrl)
		api.Get("/{file_key}", GetSignedGetUrl)
	})
	router.Route("/object-record", func(api chi.Router) {
		api.Post("/", SaveResourceInfoToRedis)
		api.Get("/{file_key}", GetResourceInfoFromRedis)
	})
	router.Route("/url", func(api chi.Router) {
		api.Post("/", SaveUrlToRedis)
	})
}
