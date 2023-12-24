package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AdminGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AdminGet called")
	w.WriteHeader(http.StatusOK)
}

// For Localstack, since the unpaid version does not support persistence through restarts
func AdminDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AdminDelete called")
	action := chi.URLParam(r, "action")
	switch action {
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
