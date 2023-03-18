package models

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetDB() *sqlx.DB {
	var db *sqlx.DB

	viper.AutomaticEnv()
	db, err := sqlx.Connect("postgres", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Error("Couldn't connect to postgresql: ", err, "\"", viper.GetString("DATABASE_URL"))
	}
	return db
}
