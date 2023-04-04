package web

import (
	"html/template"
	"ktn-go/config"
	models "ktn-go/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

func GetFeed(c *gin.Context) {
	if strings.HasSuffix(c.Param("ref"), ".xml") {
		GetFeedXml(c)
	} else {
		GetFeedHtml(c)
	}
}

func GetFeedHtml(c *gin.Context) {
	feed := models.ORMFeed{}
	ref := strings.Split(c.Param("ref"), ".")[0]
	if err := feed.GetRef(ref); err != nil {
		log.Info("404: ", err)
		c.Render(http.StatusNotFound, render.Data{})
		return
	}

	tmpl, err := models.RenderTemplate(map[string]interface{}{
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

func GetFeedXml(c *gin.Context) {
	feed := models.ORMFeed{}
	ref := strings.Split(c.Param("ref"), ".")[0]
	if err := feed.GetRef(ref); err != nil {
		log.Info("404: ", err)
		c.Render(http.StatusNotFound, render.Data{})
		return
	}

	entries := []models.ORMEntry{}
	config.DB.Select(
		&entries, "select * from entries where reference = $1 order by id desc", feed.Reference)

	feedUpdated, err := models.ConvertToRFC3339(feed.UpdatedAtRfc3339())

	feedTmpl := map[string]interface{}{
		"Reference":        feed.Reference,
		"Title":            feed.Title,
		"UpdatedAtRfc3339": feedUpdated,
		"Entries":          entries,
	}

	tmpl, err := models.RenderTemplate(feedTmpl, []string{
		"web/templates/atom.xml"})

	if err != nil {
		log.Error("KTN Web: Couldn't execute feed template: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "application/atom+xml; charset=utf-8",
		Data:        tmpl,
	})

}
