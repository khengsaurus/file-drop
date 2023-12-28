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

func GetResourceValue(fileName, fileKey, fileUrl string) string {
	return fmt.Sprintf(
		"%s%s%s%s%s",
		fileName,
		consts.RedisValDelim,
		fileKey,
		consts.RedisValDelim,
		fileUrl,
	)
}

func ParseResourceInfo(resourceVal string) (*types.ResourceInfo, error) {
	fileName := ""
	fileKey := ""
	url := ""

	if strings.Count(resourceVal, consts.RedisValDelim) >= 2 {
		// resourceVal is in the format name___key___url
		lastIndex := strings.LastIndex(resourceVal, consts.RedisValDelim)
		url = resourceVal[lastIndex+len(consts.RedisValDelim):]
		resourceVal = resourceVal[:lastIndex] // now name___key
		lastIndex = strings.LastIndex(resourceVal, consts.RedisValDelim)
		fileKey = resourceVal[lastIndex+len(consts.RedisValDelim):]
		fileName = resourceVal[:lastIndex]
	} else {
		return nil, fmt.Errorf("faield to parse resource information")
	}

	return &types.ResourceInfo{FileName: fileName, Key: fileKey, Url: url}, nil
}

func GetResourceInfoFromCtx(ctx context.Context, key string) (*types.ResourceInfo, error) {
	var resourceInfo *types.ResourceInfo
	redisCache, _ := ctx.Value(consts.RedisCacheKey).(LruCache)
	resourceVal := redisCache.Get(key)

	if resourceVal == "" {
		redisClient, err := database.GetRedisClient(ctx)
		if err != nil {
			return nil, err
		}

		resourceVal = redisClient.GetRedisValue(ctx, key)
		redisCache.Put(key, resourceVal)
	}

	resourceInfo, err := ParseResourceInfo(resourceVal)
	if err != nil {
		return nil, err
	}

	return resourceInfo, nil
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
