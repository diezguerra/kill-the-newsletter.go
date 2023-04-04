package web

import (
	"errors"
	"fmt"
	models "ktn-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func GetIndex(c *gin.Context) {
	tmpl, err := models.RenderTemplate(map[string]interface{}{}, []string{
		"web/templates/base.tmpl",
		"web/templates/index.tmpl"})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, render.Data{
		ContentType: "text/html; charset=utf-8",
		Data:        tmpl,
	})
}

func PostIndex(c *gin.Context) {
	title := c.PostForm("title")
	if len(title) == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid title"))
	}

	feed := models.ORMFeed{
		Title: title,
	}

	ref, err := feed.Save()

	if len(ref) > 0 {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/feeds/%s.html", ref))
	} else {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
