package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/gql-server/dataloader"
	"gorm.io/gorm"
)

func registerMiddleware(router *gin.RouterGroup, db *gorm.DB) {
	router.Use(dataloader.Middleware(db))
	router.Use(auth.Middleware())
}
