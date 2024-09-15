package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	ChiMiddleware "github.com/go-chi/chi/v5/middleware"
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

	fmt.Println("init Stream service")
	redisCache := utils.NewLruCache(1000, 10*time.Minute, 20*time.Minute)
	redisClient := database.InitRedisClient()
	s3Client := database.InitS3Client()

	router := chi.NewRouter()
	router.Use(middlewares.EnableCors)
	router.Use(ChiMiddleware.Timeout(30 * time.Second))
	router.Use(middlewares.WithContext(consts.RedisCacheKey, redisCache))
	router.Use(middlewares.WithContext(consts.RedisClientKey, redisClient))
	router.Use(middlewares.WithContext(consts.S3ClientKey, s3Client))

	router.HandleFunc("/ping", controllers.Ping)

	router.Route("/stream", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			rateLimiter := utils.NewRateLimiter(20, 2*time.Minute, 3*time.Minute, true)
			r.Use(rateLimiter.Handle)
			r.Get("/{file_key}", controllers.StreamResourceForView)
		})
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	if err != nil {
		panic(err)
	}
}
