package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
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

	utils.Json200(&Payload{Url: url, FileName: fileName}, w)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateRecord called")

	var req Payload
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

	nameAndGetUrl := fmt.Sprintf("%s%s%s", req.FileName, consts.RedisValDelim, getUrl)
	err = redisClient.SetRedisValue(ctx, shortestKey, nameAndGetUrl)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.Json200(&Payload{Key: shortestKey}, w)
}
