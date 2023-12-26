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

func GetRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetRecord called")
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

	concatString := redisClient.GetRedisValue(ctx, key)
	resourceInfo, err := utils.GetResourceInfo(concatString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(resourceInfo, w)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateRecord called")

	var req types.ResourceInfo
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
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

	shortestKey := redisClient.GetShortestNewKey(ctx, req.Key)
	getUrl, err := database.GetSignedGetUrl(ctx, req.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceValue := utils.GetResourceValue(req.FileName, req.Key, getUrl)
	err = redisClient.SetRedisValue(ctx, shortestKey, resourceValue)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&types.ResourceInfo{Key: shortestKey}, w)
}
