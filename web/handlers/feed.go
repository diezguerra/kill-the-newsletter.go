package web

import (
	"bytes"
	"html/template"
	models "ktn-go/models"
	"net/http"
	"strings"
	ttemplate "text/template"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

func GetFeed(c *gin.Context) {
	if strings.HasSuffix(c.Param("ref"), ".xml") {
		GetFeedXML(c)
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

func GetFeedXML(c *gin.Context) {
	feed := models.ORMFeed{}
	ref := strings.Split(c.Param("ref"), ".")[0]
	if err := feed.GetRef(ref); err != nil {
		log.Info("404: ", err)
		c.Render(http.StatusNotFound, render.Data{})
		return
	}

	atomFeed, err := ttemplate.ParseFiles("web/templates/atom.xml")
	if err != nil {
		log.Error("KTN Web: Couldn't load Atom feed template: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	executedFeed := new(bytes.Buffer)

	feedTmpl := models.Feed{
		WebUrl:      "ktngo.com",
		Reference:   "asdf",
		Title:       "Asdf",
		EmailDomain: "ktngo.com",
		Entries:     nil,
	}

	if err := atomFeed.Execute(executedFeed, feedTmpl); err != nil {
		log.Error("KTN Web: Couldn't execute feed template: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "application/atom+xml; charset=utf-8",
		Data:        executedFeed.Bytes(),
	})

}
