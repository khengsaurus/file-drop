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
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> GetRecord %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resourceInfo, err := utils.RetrieveRedisValue(r.Context(), key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(resourceInfo, w)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> CreateRecord")

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
