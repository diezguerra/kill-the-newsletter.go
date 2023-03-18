package web

import (
	models "ktn-go/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

func GetFeedHtml(c *gin.Context) {
	entry := models.ORMEntry{}
	ref := strings.Split(c.Param("ref"), ".")[0]
	if err := entry.GetRef(ref); err != nil {
		log.Info("404", err)
	}
	log.Info("Entry for ref ", ref, entry)

	tmpl, err := models.RenderTemplate(map[string]interface{}{
		"Reference":   ref,
		"WebUrl":      "ktnrs.com",
		"EmailDomain": "ktnrs.com",
		"Entry":       entry,
	}, []string{
		"../templates/base.tmpl",
		"../templates/created.tmpl"})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "text/html; charset=utf-8",
		Data:        tmpl,
	})

}
