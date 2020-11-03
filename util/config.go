package util

type ServerConfig struct {
	Version           string
	Environment       string
	AllowedOrigins    []string
	APIPath           string
	PlaygroundPath    string
	PlaygroundAPIPath string
	Port              string
	GraphQL           GqlConfig
	JWT               JWTConfig
	Hash              HashConfig
	Database          DatabaseConfig
	Redis             RedisConfig
}

type JWTConfig struct {
	Secret string
}

type HashConfig struct {
	Salt      string
	MinLength int
}

type GqlConfig struct {
	ComplexityLimit int
}

type DatabaseConfig struct {
	Host           string
	User           string
	Password       string
	Database       string
	Port           string
	MaxConnections int
}

type RedisConfig struct {
	Address string
}
