package models

// # This model works on top of the generated `feeds` SQL table
//
// ```sql
//     CREATE TABLE "feeds" (
//       "id" INTEGER PRIMARY KEY AUTOINCREMENT,
//       "createdAt" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
//       "updatedAt" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
//       "reference" TEXT NOT NULL UNIQUE,
//       "title" TEXT NOT NULL
//     );
// ```

import (
	"errors"
	"fmt"
	"ktn-go/config"
	"math/rand"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Feed struct {
	WebUrl      string
	Reference   string
	Title       string
	EmailDomain string
	Entries     []Entry
}

type NewFeed struct {
	Title     string
	Reference string
}

type ORMFeed struct {
	ID        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	Reference string `db:"reference" json:"reference"`
	Title     string `db:"title" json:"title"`
}

func (e *ORMFeed) GetRef(ref string) error {
	err := config.DB.Get(e, "select * from feeds where reference = $1 LIMIT 1", ref)

	if err != nil {
		log.Info("Failed to find Feed in db with Reference ", ref)
	}

	return err
}

// Returns a `Feed`'s `title` given its `reference`.
func (f *ORMFeed) GetTitleGivenReference(reference string, db *sqlx.DB) (string, error) {
	var title string
	err := db.Get(&title, "SELECT title FROM feeds WHERE reference = $1", reference)
	if err != nil {
		return "", err
	}
	return title, nil
}

// Checks whether a `Feed` exists given its `reference`.
func (f *ORMFeed) FeedExists(reference string, db *sqlx.DB) (bool, error) {
	var count int64
	err := db.Get(&count, "SELECT count(id) FROM feeds WHERE reference = $1", reference)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (f *ORMFeed) UpdatedAtRfc3339() string {
	entry := ORMEntry{}
	err := config.DB.Get(
		&entry,
		"select * from entries where reference = $1 ORDER BY id DESC LIMIT 1",
		f.Reference)

	if err != nil {
		log.Info("Failed to find Entries in db with Reference ", f.Reference)
		return ""
	}

	return entry.CreatedAt
}

func (f *ORMFeed) NewReference() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	if len(f.Reference) == 0 {
		f.Reference = string(b)
	}

	return string(b)
}

func (f *ORMFeed) SentinelEntry() string {
	content := map[string]interface{}{
		"Reference": f.Reference,
		"Title":     f.Title,
	}

	rendered, err := RenderTemplate(content, []string{"web/templates/sentinel_entry.tmpl"})
	if err != nil {
		log.Error("KTemplate: Couldn't render sentinel template: ", err)
		return ""
	}

	return string(rendered)
}

func (f *ORMFeed) Save() (string, error) {
	if len(f.Reference) <= 0 {
		f.Reference = f.NewReference()
	}

	tx, err := config.DB.Beginx()
	if err != nil {
		return "", err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err := tx.Exec("INSERT INTO feeds (reference, title) VALUES ($1, $2)", f.Reference, f.Title)
	if err != nil {
		return "", err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if count != 1 {
		return "", errors.New("could not create new feed")
	}

	rendered := f.SentinelEntry()
	if len(rendered) == 0 {
		return "", errors.New("couldn't render sentinel entry")
	}

	entryTitle := fmt.Sprintf("%s inbox created!", f.Title)
	_, err = tx.Exec("INSERT INTO entries (reference, title, author, content) VALUES ($1, $2, $3, $4)",
		f.Reference, entryTitle, "Kill The Newsletter", rendered)
	if err != nil {
		return "", err
	}

	return f.Reference, nil
}
