package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func GetResourceInfoFromRedis(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> GetResourceInfoFromRedis %s\n", key)
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resourceInfo, err := ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000")
	Json200(resourceInfo, w)
}

func SaveResourceInfoToRedis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> SaveResourceInfoToRedis")

	var p ResourceInfo
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
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

	shortestKey := redisClient.GetShortestNewKey(ctx, p.Key)
	getUrl, err := database.GetSignedGetUrl(ctx, p.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceValue := BuildRedisValue(p.FileName, p.Key, getUrl)
	err = redisClient.SetValue(ctx, shortestKey, string(resourceValue), 24*time.Hour)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&ResourceInfo{Key: shortestKey}, w)
}

func SaveUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> SaveUrl")

	var p UrlInfo
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	mySqlClient, err := database.GetMySqlClient(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortestKey, err := mySqlClient.GetShortestNewKey(ctx, utils.RandString(8))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = mySqlClient.WriteUrlRecord(ctx, shortestKey, p.Url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = SaveUrlToRedis(ctx, shortestKey, p.Url)
	if err != nil {
		fmt.Println(err)
	}
	Json200(&UrlInfo{Url: p.Url, Key: shortestKey}, w)
}

// --------------- types ---------------

type UrlInfo struct {
	Url string `json:"url"`
	Key string `json:"key"`
}
