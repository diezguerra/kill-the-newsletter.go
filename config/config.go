package config

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DATABASE_URL string
var DB *sqlx.DB

func LoadConfig() {
	viper.AutomaticEnv()
	log.Info("DB ", viper.GetString("DATABASE_URL"))

	DATABASE_URL = viper.GetString("DATABASE_URL")
	DB = sqlx.MustConnect("postgres", DATABASE_URL)

	viper.SetEnvPrefix("ktn")
}
