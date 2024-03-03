package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/utils"
)

func ViewFile(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> ViewFile %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redisValue, err := utils.GetRedisValue(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		utils.Redirect404(w, r)
		return
	}

	resourceInfo, err := utils.ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
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
		clientDownloadUrl := ""
		if consts.Local {
			clientDownloadUrl = os.Getenv("CLIENT_DOWNLOAD_URL_DEV")
		} else {
			clientDownloadUrl = os.Getenv("CLIENT_DOWNLOAD_URL")
		}
		clientUrl := fmt.Sprintf("%s/%s", clientDownloadUrl, key)
		http.Redirect(w, r, clientUrl, http.StatusFound)
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
