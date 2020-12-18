package main

//go:generate go run github.com/kiwisheets/gqlgencs

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/config"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/server"
	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/directive"
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
