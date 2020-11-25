package server

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerMiddleware(router *gin.RouterGroup, db *gorm.DB, cfg *util.ServerConfig) {
	router.Use(dataloader.Middleware(db))
	router.Use(auth.Middleware(db, &cfg.JWT))
}
