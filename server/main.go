package server

import (
	"github.com/gin-contrib/cors"
	"github.com/kiwisheets/gql-server/config"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Run starts a new server
func Run(cfg *config.Config, db *gorm.DB) {
	router := gin.Default()

	registerMiddleware(&router.RouterGroup, db)

	// allow access from Apollo Studio in dev mode
	if cfg.Environment == "development" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{
			"https://studio.apollographql.com",
		}
		config.AllowCredentials = true
		router.RouterGroup.Use(cors.New(config))
	}

	registerRoutes(&router.RouterGroup, cfg, db)

	log.Println("Server listening @ \"/\" on " + cfg.GraphQL.Port)
	router.Run()
}
