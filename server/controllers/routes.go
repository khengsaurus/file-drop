package controllers

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/utils"
)

var ApiRouter = func(router chi.Router) {
	router.Route("/object", func(r chi.Router) {
		r.Get("/{file_key}", GetSignedGetUrl)
		r.Group(func(r chi.Router) {
			rateLimiter := utils.NewRateLimiter(5, 2*time.Minute, 3*time.Minute, false)
			r.Use(rateLimiter.Handle)
			r.Post("/", GetSignedPutUrl)
		})
	})
	router.Route("/object-record", func(r chi.Router) {
		r.Post("/", SaveResourceInfoToRedis)
		r.Get("/{file_key}", GetResourceInfoFromRedis)
	})
	router.Route("/token", func(r chi.Router) {
		r.Get("/", GetToken)
	})
	router.Route("/url", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			rateLimiter := utils.NewRateLimiter(5, 2*time.Minute, 3*time.Minute, false)
			r.Use(rateLimiter.Handle)
			r.Post("/", SaveUrlToRedis)
		})
	})
}
