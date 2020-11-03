package config

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/joho/godotenv"
	"github.com/maxtroughear/goenv"
)

// Server retrieves config from environment variables
func Server() *util.ServerConfig {
	godotenv.Load()

	return &util.ServerConfig{
		Version:           goenv.MustGet("APP_VERSION"),
		Environment:       goenv.MustGet("ENVIRONMENT"),
		AllowedOrigins:    goenv.MustGetSlice("ALLOWED_ORIGINS"),
		APIPath:           goenv.CanGet("API_PATH", "/"),
		PlaygroundPath:    goenv.CanGet("PLAYGROUND_PATH", "/graphql"),
		PlaygroundAPIPath: goenv.CanGet("PLAYGROUND_API_PATH", "/api/"),
		Port:              goenv.MustGet("PORT"),
		JWT: util.JWTConfig{
			Secret: goenv.MustGetSecretFromEnv("JWT_SECRET_KEY"),
		},
		Hash: util.HashConfig{
			Salt:      goenv.MustGetSecretFromEnv("HASH_SALT"),
			MinLength: goenv.CanGetInt32("HASH_MIN_LENGTH", 10),
		},
		GraphQL: util.GqlConfig{
			ComplexityLimit: 200,
		},
		Database: util.DatabaseConfig{
			Host:           goenv.MustGet("POSTGRES_HOST"),
			Port:           goenv.MustGet("POSTGRES_PORT"),
			User:           goenv.MustGet("POSTGRES_USER"),
			Password:       goenv.MustGetSecretFromEnv("POSTGRES_PASSWORD"),
			Database:       goenv.MustGet("POSTGRES_DB"),
			MaxConnections: goenv.CanGetInt32("POSTGRES_MAX_CONNECTIONS", 20),
		},
		Redis: util.RedisConfig{
			Address: goenv.MustGet("REDIS_ADDRESS"),
		},
	}
}
