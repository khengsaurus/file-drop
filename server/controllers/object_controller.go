package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/types"
	"github.com/khengsaurus/file-drop/server/utils"
)

type FileInfo struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func GetSignedPutUrl(w http.ResponseWriter, r *http.Request) {
	key := utils.RandString(8)
	fmt.Printf("-> GetSignedPutUrl %s\n", key)

	var p FileInfo
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil || p.Size > 2e6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := database.GetSignedPutUrl(r.Context(), key, p.Type, p.Size)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&types.ResourceInfo{Url: url, Key: key}, w)
}

func GetSignedGetUrl(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "fileKey")
	fmt.Printf("-> GetSignedGetUrl %s\n", key)

	url, err := database.GetSignedGetUrl(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&types.ResourceInfo{Url: url}, w)
}
