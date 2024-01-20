package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/utils"
)

func StreamResource(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> StreamResource %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := utils.RetrieveRedisValue(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		utils.Redirect404(w, r)
		return
	}

	resourceInfo, err := utils.ParseRedisValue(redisValue)
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
		http.Redirect(w, r, clientUrl, http.StatusTemporaryRedirect)
		return
	}

	result, err := database.GetObject(ctx, resourceInfo.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(buffer)

	http.ServeContent(w, r, resourceInfo.FileName, time.Now(), reader)
}
