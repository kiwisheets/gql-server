package server

import (
	"log"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Run starts a new server
func Run(cfg *util.ServerConfig, db *gorm.DB) {
	router := gin.Default()

	registerMiddleware(&router.RouterGroup, db, cfg)

	registerRoutes(&router.RouterGroup, cfg, db)

	log.Println("Server listening @ \"/\" on " + cfg.Port)
	router.Run()
}
