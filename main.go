package main

//go:generate go run github.com/99designs/gqlgen generate

import (
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/config"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/server"
	"github.com/emvi/hide"
)

func main() {
	cfg := config.Server()

	hide.UseHash(hide.NewHashID(cfg.Hash.Salt, cfg.Hash.MinLength))

	// Give time for the database to start
	time.Sleep(5000 * time.Millisecond)

	// connect to db
	db := orm.Init(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	server.Run(cfg, db)
}
