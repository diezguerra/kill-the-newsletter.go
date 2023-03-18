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
	CreatedAt string `db:"createdAt" json:"created_at"`
	UpdatedAt string `db:"updatedAt" json:"updated_at"`
	Reference string `db:"reference" json:"reference"`
	Title     string `db:"title" json:"title"`
	Author    string `db:"author" json:"author"`
	Contenet  string `db:"content" json:"content"`
}

func (e ORMEntry) GetRef(ref string) error {
	log.Info("Getting Entry for ref", ref)
	err := config.DB.Get(&e, "select from entries where reference = $1 LIMIT 1", ref)

	if err != nil {
		log.Info("Failed to find Entry with Reference ", ref)
	}

	return err
}

func (e ORMEntry) String() string {
	tmpl, err := RenderTemplate(e, []string{
		"../web/templates/sentinel_entry.tmpl",
	})
	if err != nil {
		return err.Error()
	}
	return string(tmpl)

}
