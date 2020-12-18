package server

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/auth"
	"gorm.io/gorm"
)

func registerMiddleware(router *gin.RouterGroup, db *gorm.DB) {
	router.Use(dataloader.Middleware(db))
	router.Use(auth.Middleware())
}
