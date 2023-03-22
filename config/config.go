package config

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DATABASE_URL string
var DB *sqlx.DB
var TemplateVariables map[string]string

func LoadConfig() {
	viper.AutomaticEnv()
	log.Info("DB ", viper.GetString("DATABASE_URL"))

	DATABASE_URL = viper.GetString("DATABASE_URL")
	DB = sqlx.MustConnect("postgres", DATABASE_URL)

	viper.SetEnvPrefix("ktn")

	TemplateVariables = make(map[string]string)

	TemplateVariables["EmailDomain"] = viper.GetString("DOMAIN")
	TemplateVariables["WebUrl"] = viper.GetString("URL")
}
