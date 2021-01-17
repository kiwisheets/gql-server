package main

//go:generate go run github.com/maxtroughear/gqlgencs

import (
	"os"

	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/directive"
	"github.com/kiwisheets/gql-server/config"
	"github.com/kiwisheets/gql-server/orm"
	"github.com/kiwisheets/gql-server/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Server()

	logrus.SetOutput(os.Stdout)
	if cfg.Environment == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	hide.UseHash(hide.NewHashID(cfg.Hash.Salt, cfg.Hash.MinLength))

	// connect to db
	db := orm.Init(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	directive.Development(cfg.GraphQL.Environment == "development")

	server.Run(cfg, db)
}
