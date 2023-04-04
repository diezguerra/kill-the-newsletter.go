package models

import (
	"ktn-go/config"

	log "github.com/sirupsen/logrus"
)

type Entry struct {
	id        int
	Title     string
	Reference string
	CreatedAt string
	UpdatedAt string
	Author    string
	Content   string
}

type ORMEntry struct {
	Id        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	Reference string `db:"reference" json:"reference"`
	Title     string `db:"title" json:"title"`
	Author    string `db:"author" json:"author"`
	Content   string `db:"content" json:"content"`
}

func (e *ORMEntry) GetRef(ref string) error {
	log.Info("Fetching Entry db record for ref", ref)
	err := config.DB.Get(e, "select * from entries where reference = $1 LIMIT 1", ref)

	if err != nil {
		log.Info("Failed to find Entry in db with Reference ", ref)
	}

	return err
}

func (e *ORMEntry) UpdatedAtRfc3339() string {
	feedUpdatedAt, _ := ConvertToRFC3339(e.CreatedAt)
	return feedUpdatedAt
}

func (e *ORMEntry) String() string {

	tmplVars := map[string]interface{}{
		"CreatedAt": e.CreatedAt,
		"Reference": e.Reference,
		"Title":     e.Title,
		"Author":    e.Author,
		"Content":   e.Content}

	tmpl, err := RenderTemplate(tmplVars, []string{
		"web/templates/sentinel_entry.tmpl",
	})

	if err != nil {
		return err.Error()
	}

	return string(tmpl)

}
