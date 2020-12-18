package main

//go:generate go run github.com/kiwisheets/gqlgencs

import (
	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/directive"
	"github.com/kiwisheets/gql-server/config"
	"github.com/kiwisheets/gql-server/orm"
	"github.com/kiwisheets/gql-server/server"
)

func main() {
	cfg := config.Server()

	hide.UseHash(hide.NewHashID(cfg.Hash.Salt, cfg.Hash.MinLength))

	// connect to db
	db := orm.Init(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	directive.Development(cfg.GraphQL.Environment == "development")

	server.Run(cfg, db)
}
