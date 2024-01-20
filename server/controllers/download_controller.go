package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func StreamResourceForDownload(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> StreamResourceForDownload %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := utils.RetrieveRedisValue(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceInfo, err := utils.ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := database.GetObject(ctx, resourceInfo.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", resourceInfo.FileName))
	w.Header().Set("Cache-Control", "no-store")

	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error streaming resource as HTTP response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
