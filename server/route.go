package server

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/handler"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func graphqlHandler(db *gorm.DB, cfg *util.ServerConfig) gin.HandlerFunc {
	gql := handler.GraphqlHandler(db, cfg)
	return func(c *gin.Context) {
		gql.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler(cfg *util.ServerConfig) gin.HandlerFunc {
	playground := playground.Handler("GraphQL playground", cfg.PlaygroundAPIPath)
	return func(c *gin.Context) {
		playground.ServeHTTP(c.Writer, c.Request)
	}
}

func registerRoutes(router *gin.RouterGroup, cfg *util.ServerConfig, db *gorm.DB) {
	// support GET for automatic persisted queries
	router.GET(cfg.APIPath, graphqlHandler(db, cfg))
	router.POST(cfg.APIPath, graphqlHandler(db, cfg))
	router.OPTIONS(cfg.APIPath, graphqlHandler(db, cfg))

	if cfg.Environment == "development" {
		router.GET(cfg.PlaygroundPath, playgroundHandler(cfg))
	}
}
