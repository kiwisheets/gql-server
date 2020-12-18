package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

// Run starts a new server
func Run(cfg *util.ServerConfig, db *gorm.DB) {
	router := gin.Default()

	registerMiddleware(&router.RouterGroup, db)

	registerRoutes(&router.RouterGroup, cfg, db)

	log.Println("Server listening @ \"/\" on " + cfg.GraphQL.Port)
	router.Run()
}
