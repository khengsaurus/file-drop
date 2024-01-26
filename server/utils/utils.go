package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/types"
)

func BuildRedisValue(fileName, fileKey, fileUrl string) consts.RedisResourceValue {
	return consts.RedisResourceValue(fmt.Sprintf(
		"%s%s%s%s%d%s%s",
		fileName,
		consts.RedisValDelim,
		fileKey,
		consts.RedisValDelim,
		time.Now().Unix(),
		consts.RedisValDelim,
		fileUrl,
	))
}

func ParseRedisValue(resourceValue consts.RedisResourceValue) (*types.ResourceInfo, error) {
	fullString := string(resourceValue)
	values := strings.Split(fullString, consts.RedisValDelim)

	if len(values) < 4 {
		return nil, fmt.Errorf("failed to parse resource information")
	}

	uploadedAt, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource information")
	}

	return &types.ResourceInfo{
		FileName:   values[0],
		Key:        values[1],
		UploadedAt: uploadedAt,
		Url:        values[3],
	}, nil
}

func RetrieveRedisValue(ctx context.Context, key string) (consts.RedisResourceValue, error) {
	redisCache, _ := ctx.Value(consts.RedisCacheKey).(LruCache)
	redisValue := redisCache.Get(key)

	if redisValue == "" {
		redisClient, err := database.GetRedisClient(ctx)
		if err != nil {
			return "", err
		}

		redisValue = redisClient.RetrieveRedisValue(ctx, key)
		if redisValue != "" {
			redisCache.Put(key, redisValue)
		}
	}

	return consts.RedisResourceValue(redisValue), nil
}

func Json200(payload any, w http.ResponseWriter) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	url := ""
	if consts.Local {
		url = os.Getenv("CLIENT_HOME_URL_DEV")
	} else {
		url = os.Getenv("CLIENT_HOME_URL")
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Redirect404(w http.ResponseWriter, r *http.Request) {
	client404 := ""
	if consts.Local {
		client404 = os.Getenv("CLIENT_404_URL_DEV")
	} else {
		client404 = os.Getenv("CLIENT_404_URL")
	}
	http.Redirect(w, r, client404, http.StatusTemporaryRedirect)
}

func WriteImageHTML(title string, src string, w http.ResponseWriter) error {
	tmpl, err := template.New("imagePage").Parse(ImagePageHtml)
	if err != nil {
		return err
	}

	page := types.HtmlPageImg{Title: title, Src: src}

	return tmpl.Execute(w, page)
}
