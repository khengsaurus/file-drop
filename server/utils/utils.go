package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/types"
)

func BuildRedisValue(fileName, fileKey, fileUrl string) consts.RedisResourceValue {
	return consts.RedisResourceValue(fmt.Sprintf(
		"%s%s%s%s%s",
		fileName,
		consts.RedisValDelim,
		fileKey,
		consts.RedisValDelim,
		fileUrl,
	))
}

func ParseRedisValue(resourceValue consts.RedisResourceValue) (*types.ResourceInfo, error) {
	fileName, fileKey, url, val := "", "", "", string(resourceValue)

	if strings.Count(val, consts.RedisValDelim) >= 2 {
		lastIndex := strings.LastIndex(val, consts.RedisValDelim)
		url = val[lastIndex+len(consts.RedisValDelim):]
		val = val[:lastIndex] // now name___key
		lastIndex = strings.LastIndex(val, consts.RedisValDelim)
		fileKey = val[lastIndex+len(consts.RedisValDelim):]
		fileName = val[:lastIndex]
	} else {
		return nil, fmt.Errorf("faield to parse resource information")
	}

	return &types.ResourceInfo{FileName: fileName, Key: fileKey, Url: url}, nil
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

func WriteImageHTML(title string, src string, w http.ResponseWriter) error {
	tmpl, err := template.New("imagePage").Parse(ImagePageHtml)
	if err != nil {
		return err
	}

	page := types.HtmlPageImg{Title: title, Src: src}

	return tmpl.Execute(w, page)
}
