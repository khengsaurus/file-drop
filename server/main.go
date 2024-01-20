package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ChiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/controllers"
	"github.com/khengsaurus/file-drop/server/database"
	"github.com/khengsaurus/file-drop/server/middlewares"
	"github.com/khengsaurus/file-drop/server/utils"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}

	redisCache := utils.LruCacheConstructor(20, 5*time.Minute, 10*time.Minute)
	redisClient := database.InitRedisClient()
	s3Client := database.InitS3Client()

	router := chi.NewRouter()
	router.Use(middlewares.EnableCors)
	router.Use(ChiMiddleware.Timeout(30 * time.Second))
	router.Use(middlewares.WithContext(consts.RedisCacheKey, redisCache))
	router.Use(middlewares.WithContext(consts.RedisClientKey, redisClient))
	router.Use(middlewares.WithContext(consts.S3ClientKey, s3Client))

	router.Route("/api", func(r chi.Router) {
		controllers.ApiRouter(r)
	})
	router.Route("/file", func(r chi.Router) {
		r.Get("/{file_key}", controllers.ViewFile)
	})
	router.Route("/stream", func(r chi.Router) {
		r.Get("/{file_key}", controllers.StreamResource)
	})
	router.Route("/download", func(r chi.Router) {
		r.Get("/{file_key}", controllers.StreamResourceForDownload)
	})
	router.Route("/url", func(r chi.Router) {
		r.Get("/{url_key}", controllers.ReditectToUrlFromRedis)
	})

	// Dev
	if consts.Local {
		router.HandleFunc("/test", test)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	if err != nil {
		panic(err)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
