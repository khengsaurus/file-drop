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

	fmt.Println("init URL service")
	mySqlClient := database.InitMySqlConnection(24*time.Hour, time.Hour)
	redisCache := utils.NewLruCache(1000, 10*time.Minute, 20*time.Minute)
	redisClient := database.InitRedisClient()

	router := chi.NewRouter()
	router.Use(middlewares.EnableCors)
	router.Use(ChiMiddleware.Timeout(30 * time.Second))
	router.Use(middlewares.WithContext(consts.MySqlClientKey, mySqlClient))
	router.Use(middlewares.WithContext(consts.RedisCacheKey, redisCache))
	router.Use(middlewares.WithContext(consts.RedisClientKey, redisClient))

	router.HandleFunc("/ping", ping)

	router.Route("/url", func(r chi.Router) {
		r.Get("/{url_key}", controllers.RedirectToUrl)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
