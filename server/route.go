package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/kiwisheets/gql-server/config"
	"github.com/kiwisheets/gql-server/handler"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

func graphqlHandler(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	gql := handler.GraphqlHandler(db, cfg)
	return func(c *gin.Context) {
		gql.ServeHTTP(c.Writer, c.Request)
	}
}

func healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Healthy")
	}
}

func playgroundHandler(cfg *util.GqlConfig) gin.HandlerFunc {
	playground := playground.Handler("GraphQL playground", cfg.PlaygroundAPIPath)
	return func(c *gin.Context) {
		playground.ServeHTTP(c.Writer, c.Request)
	}
}

func registerRoutes(router *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	router.GET("/health", healthHandler())

	// support GET for automatic persisted queries
	router.GET(cfg.GraphQL.APIPath, graphqlHandler(db, cfg))
	router.POST(cfg.GraphQL.APIPath, graphqlHandler(db, cfg))

	if cfg.Environment == "development" {
		router.GET(cfg.GraphQL.PlaygroundPath, playgroundHandler(&cfg.GraphQL))
	}
}
