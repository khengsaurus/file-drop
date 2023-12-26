package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/types"
	"github.com/khengsaurus/file-drop/server/utils"
)

func GetSignedPutUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSignedPutUrl called")

	key := uuid.New().String()
	url, err := database.GetSignedPutUrl(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&types.ResourceInfo{Url: url, Key: key}, w)
}

func GetSignedGetUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSignedGetUrl called")
	key := chi.URLParam(r, "fileKey")

	url, err := database.GetSignedGetUrl(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&types.ResourceInfo{Url: url}, w)
}
