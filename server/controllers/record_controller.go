package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/types"
	"github.com/khengsaurus/file-drop/server/utils"
)

func GetResourceInfoFromRedis(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> GetResourceInfoFromRedis %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redisValue, err := utils.RetrieveRedisValue(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resourceInfo, err := utils.ParseRedisValue(redisValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000")
	utils.Json200(resourceInfo, w)
}

func SaveResourceInfoToRedis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> SaveResourceInfoToRedis")

	var p types.ResourceInfo
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

	resourceValue := utils.BuildRedisValue(p.FileName, p.Key, getUrl)
	err = redisClient.SetRedisValue(ctx, shortestKey, string(resourceValue))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&types.ResourceInfo{Key: shortestKey}, w)
}

func SaveUrlToRedis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> SaveUrlToRedis")

	var p types.UrlInfo
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

	key := uuid.New().String()
	shortestKey := redisClient.GetShortestNewKey(ctx, key)
	err = redisClient.SetRedisValue(ctx, shortestKey, p.Url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&types.UrlInfo{Url: p.Url, Key: shortestKey}, w)
}
