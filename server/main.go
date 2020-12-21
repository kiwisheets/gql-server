package server

import (
	"github.com/kiwisheets/gql-server/config"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Run starts a new server
func Run(cfg *config.Config, db *gorm.DB) {
	router := gin.Default()

	registerMiddleware(&router.RouterGroup, db)

	registerRoutes(&router.RouterGroup, cfg, db)

	log.Println("Server listening @ \"/\" on " + cfg.GraphQL.Port)
	router.Run()
}
