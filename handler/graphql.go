package handler

import (
	"fmt"
	"net/http"
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/directive"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/resolver"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

// GraphqlHandler constructs and returns a http handler
func GraphqlHandler(db *gorm.DB, cfg *util.ServerConfig) http.Handler {
	c := generated.Config{
		Resolvers: &resolver.Resolver{
			DB:  db,
			Cfg: cfg,
		},
		Directives: directive.Register(),
	}

	// init APQ cache
	cache, err := newCache(cfg.Redis.Address, 24*time.Hour)
	if err != nil {
		panic(fmt.Errorf("cannot create APQ cache: %v", err))
	}

	gqlHandler := handler.New(generated.NewExecutableSchema(c))

	gqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.MultipartForm{})

	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: cache,
	})

	gqlHandler.Use(extension.Introspection{})

	// gqlHandler.Use(&extension.ComplexityLimit{
	// 	Func: func(ctx context.Context, rc *graphql.OperationContext) int {
	// 		return cfg.GraphQL.ComplexityLimit
	// 	},
	// })

	return gqlHandler
}
