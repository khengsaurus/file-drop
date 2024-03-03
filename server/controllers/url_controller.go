package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func RedirectToUrl(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "url_key")
	fmt.Printf("-> RedirectToUrl %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := utils.GetRedisValue(ctx, key)
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

	if urlEntry != nil {
		var url = ""
		if strings.HasPrefix(urlEntry.Link, "http") {
			url = urlEntry.Link
		} else {
			url = fmt.Sprintf("https://%s", urlEntry.Link)
		}
		http.Redirect(w, r, url, http.StatusFound)
		utils.SaveUrlToRedis(ctx, urlEntry.Id, urlEntry.Link)
	} else {
		utils.RedirectHome(w, r)
	}

}
