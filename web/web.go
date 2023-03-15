package ktn

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpServer() *http.Server {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return srv
}
