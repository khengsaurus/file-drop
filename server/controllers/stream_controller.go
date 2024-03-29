package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
)

func StreamResource(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "file_key")
	fmt.Printf("-> StreamResource %s\n", key)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	redisValue, err := GetRedisValue(ctx, key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(redisValue) == "" {
		Redirect404(w, r)
		return
	}

	resourceInfo, err := ParseRedisValue(redisValue)
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
		http.Redirect(w, r, clientUrl, http.StatusFound)
		return
	}

	result, err := database.GetObject(ctx, resourceInfo.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(buffer)

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", resourceInfo.FileName))
	http.ServeContent(w, r, resourceInfo.FileName, time.Unix(resourceInfo.UploadedAt, 0), reader)
}
