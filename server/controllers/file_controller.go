package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
)

type fileInfo struct {
	Size int    `json:"size"`
	Type string `json:"type"`
}

func GetSignedPutUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> GetSignedPutUrl")

	clientToken, _ := r.Cookie(consts.ClientCookieName)
	if clientToken == nil || clientToken.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var p fileInfo
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil || p.Size > 2e6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key, url, err := database.GetSignedPutUrl(r.Context(), p.Type, p.Size, 0, clientToken.Value)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&ResourceInfo{Url: url, Key: key}, w)
}

func GetSignedGetUrl(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> GetSignedGetUrl %s\n", key)

	url, err := database.GetSignedGetUrl(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&ResourceInfo{Url: url}, w)
}

func ViewFile(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> ViewFile %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redisValue, err := GetRedisValue(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		Redirect404(w, r)
		return
	}

	resourceInfo, err := ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	extension := filepath.Ext(resourceInfo.FileName)
	switch strings.ToLower(extension) {
	case ".img", ".jpeg", ".jpg", ".png", ".svg":
		err = WriteImageHTML(resourceInfo.FileName, resourceInfo.Url, w)
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

func StreamResourceForDownload(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> StreamResourceForDownload %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := GetRedisValue(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		Redirect404(w, r)
		return
	}

	resourceInfo, err := ParseRedisValue(redisValue)
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

func StreamResourceForView(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> StreamResourceForView %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := GetRedisValue(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		Redirect404(w, r)
		return
	}

	resourceInfo, err := ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resourceInfo == nil {
		clientUrl := ""
		if consts.Local {
			clientUrl = os.Getenv("CLIENT_BASE_URL_DEV")
		} else {
			clientUrl = os.Getenv("CLIENT_BASE_URL")
		}
		http.Redirect(w, r, clientUrl, http.StatusFound)
		return
	}

	result, err := database.GetObject(ctx, resourceInfo.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(buffer)

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", resourceInfo.FileName))
	http.ServeContent(w, r, resourceInfo.FileName, time.Unix(resourceInfo.UploadedAt, 0), reader)
}
