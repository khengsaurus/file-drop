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
)

var (
	route_api   = "/api"
	route_admin = "/admin"
	route_test  = "/test"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}

	redisClient := database.InitRedisClient()
	s3Client := database.InitS3Client()

	router := chi.NewRouter()
	router.Use(middlewares.EnableCors)
	router.Use(ChiMiddleware.Timeout(30 * time.Second))

	router.Route(route_api, func(restRouter chi.Router) {
		restRouter.Use(middlewares.WithContext(consts.S3ClientKey, s3Client))
		restRouter.Use(middlewares.WithContext(consts.RedisClientKey, redisClient))
		controllers.RestRouter(restRouter)
	})

	// Dev
	if consts.Local {
		router.HandleFunc(route_test, test)
		router.Route(route_admin, func(adminRouter chi.Router) {
			adminRouter.Use(middlewares.AdminValidation)
			controllers.AdminRouter(adminRouter)
		})
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
