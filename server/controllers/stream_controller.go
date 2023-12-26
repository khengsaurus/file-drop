package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func StreamResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("StreamResource called")
	key := chi.URLParam(r, "file_key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisClient, err := database.GetRedisClient(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceVal := redisClient.GetRedisValue(ctx, key)
	resourceInfo, err := utils.GetResourceInfo(resourceVal)
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
