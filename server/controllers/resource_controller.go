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

	nameAndGetUrl := redisClient.GetRedisValue(ctx, key)
	lastIndex := strings.LastIndex(nameAndGetUrl, consts.RedisValDelim)
	fileName := ""
	url := ""

	if lastIndex != -1 {
		fileName = nameAndGetUrl[:lastIndex]
		url = nameAndGetUrl[lastIndex+len(consts.RedisValDelim):]
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	extension := filepath.Ext(fileName)
	switch strings.ToLower(extension) {
	case ".img", ".jpeg", ".jpg", ".png", ".svg":
		err = utils.WriteImageHTML(fileName, url, w)
	case ".json", ".text", ".txt", ".html":
		http.Redirect(w, r, url, http.StatusFound)
	default:
		clientBaseUrl := os.Getenv("CLIENT_URL")
		if consts.Local {
			clientBaseUrl = os.Getenv("CLIENT_URL_DEV")
		}
		clientUrl := fmt.Sprintf("%s/file/%s", clientBaseUrl, key)
		http.Redirect(w, r, clientUrl, http.StatusFound)
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
