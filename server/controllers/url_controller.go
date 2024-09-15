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

type urlInfo struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func SaveUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> SaveUrl")

	var p urlInfo
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

	entry, err := mySqlClient.WriteUrlRecord(ctx, p.Url, utils.GetRecordExpiryRef(r), 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = SaveUrlToRedis(ctx, entry.Id, p.Url)
	if err != nil {
		fmt.Println(err)
	}
	Json200(&urlInfo{Url: p.Url, Key: entry.Id}, w)
}

func RedirectToUrl(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "url_key")
	fmt.Printf("-> RedirectToUrl %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := GetRedisValue(ctx, key)
	if err == nil && redisValue != "" {
		var url = ""
		if strings.HasPrefix(string(redisValue), "http") {
			url = string(redisValue)
		} else {
			url = fmt.Sprintf("https://%s", redisValue)
		}
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	mySqlClient := ctx.Value(consts.MySqlClientKey).(*database.MySqlClient)
	urlEntry, err := mySqlClient.GetUrlRecordById(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if urlEntry == nil {
		RedirectHome(w, r)
		return
	}

	var url = ""
	if strings.HasPrefix(urlEntry.Link, "http") {
		url = urlEntry.Link
	} else {
		url = fmt.Sprintf("https://%s", urlEntry.Link)
	}
	http.Redirect(w, r, url, http.StatusFound)
	SaveUrlToRedis(ctx, urlEntry.Id, urlEntry.Link)

}
