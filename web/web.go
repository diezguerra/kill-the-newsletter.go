package ktn

import (
	"bytes"
	handlers "ktn-go/web/handlers"
	"net/http"

	htemplate "html/template"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

func HttpServer() *http.Server {
	r := gin.Default()

	r.Static("/static", "./web/static")

	r.GET("/", func(c *gin.Context) {
		idxTmpl, err := htemplate.ParseFiles(
			"web/templates/base.tmpl",
			"web/templates/index.tmpl")
		if err != nil {
			log.Error("KTN Web: Couldn't load Index template: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		executedIdx := new(bytes.Buffer)
		if err := idxTmpl.Execute(executedIdx, struct{ WebUrl string }{WebUrl: "ktngo.com"}); err != nil {
			log.Error("KTN Web: Couldn't execute index template: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.Render(http.StatusOK, render.Data{
			ContentType: "text/html; charset=utf-8",
			Data:        executedIdx.Bytes(),
		})
	})

	r.GET("/feeds/asdf.xml", func(c *gin.Context) {

	})

	r.GET("/feeds/:ref", handlers.GetFeed)

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	return srv
}
