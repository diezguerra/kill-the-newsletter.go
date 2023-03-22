package web

import (
	"html/template"
	models "ktn-go/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

func GetFeedHtml(c *gin.Context) {
	feed := models.ORMFeed{}
	ref := strings.Split(c.Param("ref"), ".")[0]
	if err := feed.GetRef(ref); err != nil {
		log.Info("404: ", err)
		c.Render(http.StatusNotFound, render.Data{})
		return
	}

	tmpl, err := models.RenderHTMLTemplate(map[string]interface{}{
		"Reference": ref,
		"Entry":     template.HTML(feed.SentinelEntry()),
	}, []string{
		"web/templates/base.tmpl",
		"web/templates/created.tmpl"})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "text/html; charset=utf-8",
		Data:        tmpl,
	})

}
