package config

import (
	"log"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util/key"
	"github.com/joho/godotenv"
	"github.com/kiwisheets/util"
	"github.com/maxtroughear/goenv"
)

// Server retrieves config from environment variables
func Server() *util.ServerConfig {
	godotenv.Load()

	jwtPrivateKey, err := key.ParseEcPrivateKeyFromPemStr(goenv.MustGetSecretFromEnv("JWT_EC_PRIVATE_KEY"))
	if err != nil {
		log.Println("failed to parse JWT_EC_PRIVATE_KEY")
		log.Fatal(err)
	}

	jwtPublicKey, err := key.ParseEcPublicKeyFromPemStr(goenv.MustGetSecretFromEnv("JWT_EC_PUBLIC_KEY"))
	if err != nil {
		log.Println("failed to parse JWT_EC_PUBLIC_KEY")
		log.Fatal(err)
	}

	return &util.ServerConfig{
		Version:     goenv.MustGet("APP_VERSION"),
		Environment: goenv.MustGet("ENVIRONMENT"),
		JWT: util.JWTConfig{
			PrivateKey: jwtPrivateKey,
			PublicKey:  jwtPublicKey,
		},
		Hash: util.HashConfig{
			Salt:      goenv.MustGetSecretFromEnv("HASH_SALT"),
			MinLength: goenv.CanGetInt32("HASH_MIN_LENGTH", 10),
		},
		GraphQL: util.GqlConfig{
			APIPath:           goenv.CanGet("API_PATH", "/"),
			ComplexityLimit:   200,
			PlaygroundPath:    goenv.CanGet("PLAYGROUND_PATH", "/graphql"),
			PlaygroundAPIPath: goenv.CanGet("PLAYGROUND_API_PATH", "/api/"),
			Port:              goenv.MustGet("PORT"),
			Environment:       goenv.MustGet("ENVIRONMENT"),
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
