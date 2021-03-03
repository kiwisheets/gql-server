package main

//go:generate go run github.com/maxtroughear/gqlgencs

import (
	"github.com/kiwisheets/auth/directive"
	"github.com/kiwisheets/gql-server/config"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/model"
	"github.com/kiwisheets/gql-server/resolver"
	"github.com/kiwisheets/server/graphqlapi"
)

func main() {
	cfg := config.Server()

	app := graphqlapi.NewDefault()
	defer app.Shutdown()

	// connect to db
	db := model.Init(&cfg.Database)
	defer db.Close()

	directive.Development(cfg.GraphQL.Environment == "development")

	c := generated.Config{
		Resolvers: &resolver.Resolver{
			DB:  db.DB,
			Cfg: cfg,
		},
		Directives: generated.DirectiveRoot{
			IsAuthenticated:       directive.IsAuthenticated,
			IsSecureAuthenticated: directive.IsSecureAuthenticated,
			HasPerm:               directive.HasPerm,
			HasPerms:              directive.HasPerms,
		},
	}

	server := app.SetupServer(generated.NewExecutableSchema(c), &cfg.GraphQL, db.DB)
	server.RouterGroup.Use(dataloader.Middleware(db.DB))
	server.Run(app.Logger)
}
