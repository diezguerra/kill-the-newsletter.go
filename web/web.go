package ktn

import (
	handlers "ktn-go/web/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpServer() *http.Server {
	r := gin.Default()

	r.Static("/static", "./web/static")

	r.GET("/", handlers.GetIndex)
	r.GET("/feeds/:ref", handlers.GetFeed)

	r.POST("/", handlers.PostIndex)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	return srv
}
