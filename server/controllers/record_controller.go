package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
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

	clientToken, _ := r.Cookie(consts.ClientCookieName)
	if clientToken == nil || clientToken.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	placeholderValue := redisClient.GetValue(ctx, p.Key)
	expectedPlaceholder := fmt.Sprintf("%s_%s", consts.RedisValPlaceholderPrefix, clientToken.Value)
	if placeholderValue != expectedPlaceholder {
		fmt.Println(fmt.Errorf("write to invalid redis-key"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	getUrl, err := database.GetSignedGetUrl(ctx, p.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceValue := BuildRedisValue(p.FileName, p.Key, getUrl)
	err = redisClient.SetValue(ctx, p.Key, string(resourceValue), 24*time.Hour)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Json200(&ResourceInfo{Key: p.Key}, w)
}
