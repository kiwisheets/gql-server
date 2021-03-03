package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/kiwisheets/gql-server/util/key"
	"github.com/kiwisheets/util"
	"github.com/maxtroughear/goenv"
	"github.com/sirupsen/logrus"
)

const (
	jwtEcPrivateKey = "JWT_EC_PRIVATE_KEY_FILE"
)

var jwtEcPrivateKeyFilename string

var config Config

// Server retrieves config from environment variables
func Server() *Config {
	godotenv.Load()

	jwtEcPrivateKeyFilename = goenv.MustGet(jwtEcPrivateKey)

	config = Config{
		Version:     goenv.MustGet("APP_VERSION"),
		Environment: goenv.MustGet("ENVIRONMENT"),
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
			Redis: util.RedisConfig{
				Address: goenv.MustGet("REDIS_ADDRESS"),
			},
		},
		Database: util.DatabaseConfig{
			Host:           goenv.MustGet("POSTGRES_HOST"),
			Port:           goenv.MustGet("POSTGRES_PORT"),
			User:           goenv.MustGet("POSTGRES_USER"),
			Password:       goenv.MustGetSecretFromEnv("POSTGRES_PASSWORD"),
			Database:       goenv.MustGet("POSTGRES_DB"),
			MaxConnections: goenv.CanGetInt32("POSTGRES_MAX_CONNECTIONS", 20),
		},
	}

	registerFileWatchers()
	privateKey, err := key.ParseEcPrivateKeyFromPemStr(goenv.MustGetSecretFromFile(jwtEcPrivateKeyFilename))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": jwtEcPrivateKeyFilename,
		}).Fatal("failed to parse key, check it exists and if it is a valid EC256 private key")
	}
	config.JWT.mutex.Lock()
	config.JWT.privateKey = privateKey
	config.JWT.mutex.Unlock()

	return &config
}

func registerFileWatchers() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Error("failed to initialise file watcher, will load without watching")
	}
	err = watcher.Add(jwtEcPrivateKeyFilename)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": jwtEcPrivateKeyFilename,
		}).Error("failed to watch file, check it exists")
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logrus.Error("watcher error")
					continue
				}
				logrus.WithFields(logrus.Fields{
					"event": event,
				}).Debug("file watch event")
				if event.Op&fsnotify.Write == fsnotify.Write {
					reloadFile(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logrus.Error("watcher error")
					continue
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("watcher error")
			}
		}
	}()
}

func reloadFile(filename string) {
	if filename == jwtEcPrivateKeyFilename {
		keyString := goenv.MustGetSecretFromFile(jwtEcPrivateKeyFilename)
		if keyString == "" {
			return
		}

		logrus.WithFields(logrus.Fields{
			"filename": filename,
		}).Info("reloading file")

		privateKey, err := key.ParseEcPrivateKeyFromPemStr(keyString)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"filename": jwtEcPrivateKeyFilename,
			}).Error("failed to parse key, check it exists and if it is a valid EC256 private key. continuing with existing private key")
			return
		}
		config.JWT.mutex.Lock()
		config.JWT.privateKey = privateKey
		config.JWT.mutex.Unlock()
	}
}
