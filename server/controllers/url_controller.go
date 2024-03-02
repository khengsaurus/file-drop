package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/utils"
)

func ReditectToUrlFromRedis(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "url_key")
	fmt.Printf("-> ReditectToUrlFromRedis %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redisValue, err := utils.RetrieveRedisValue(r.Context(), key)
	if redisValue == "" || err != nil {
		if err != nil {
			fmt.Println(err)
		}
		utils.RedirectHome(w, r)
		return
	}

	var url = ""
	if strings.HasPrefix(string(redisValue), "http") {
		url = string(redisValue)
	} else {
		url = fmt.Sprintf("https://%s", redisValue)
	}

	http.Redirect(w, r, url, http.StatusFound)
}
