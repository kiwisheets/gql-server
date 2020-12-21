package config

import "github.com/kiwisheets/util"

type Config struct {
	Version     string
	Environment string
	GraphQL     util.GqlConfig
	JWT         JWTConfig
	Hash        util.HashConfig
	Database    util.DatabaseConfig
	Redis       util.RedisConfig
}
