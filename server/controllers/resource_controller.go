package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func ViewResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ViewResource called")
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	extension := filepath.Ext(resourceInfo.FileName)
	switch strings.ToLower(extension) {
	case ".img", ".jpeg", ".jpg", ".png", ".svg":
		err = utils.WriteImageHTML(resourceInfo.FileName, resourceInfo.Url, w)
	case ".json", ".text", ".txt", ".html":
		http.Redirect(w, r, resourceInfo.Url, http.StatusFound)
	default:
		clientBaseUrl := os.Getenv("CLIENT_FILE_URL")
		if consts.Local {
			clientBaseUrl = os.Getenv("CLIENT_FILE_URL_DEV")
		}
		clientUrl := fmt.Sprintf("%s/%s", clientBaseUrl, key)
		http.Redirect(w, r, clientUrl, http.StatusFound)
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
