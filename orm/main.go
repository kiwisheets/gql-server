package orm

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kiwisheets/auth/seed"
	"github.com/kiwisheets/gql-server/config"
	"github.com/kiwisheets/gql-server/orm/migration"
	internalseed "github.com/kiwisheets/gql-server/orm/seed"
	"github.com/kiwisheets/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init connects to and initialises the database
func Init(cfg *config.Config) *gorm.DB {
	dbCfg := cfg.Database

	connectionString := constructConnectionString(&dbCfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectionString,
	}), &gorm.Config{
		AllowGlobalUpdate: false,
		Logger:            logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Println("Failed to connect to db")
		log.Println(dbCfg.Host)
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Println("Failed to connect to db")
		log.Println(dbCfg.Host)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(dbCfg.MaxConnections)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)

	if cfg.Environment == "development" {
		// clear db
		// note: does not drop tables used for many2many relationships, please bare this in mind!
		migration.DropAll(db)
		db.Config.Logger = logger.Default.LogMode(logger.Info)
	}

	migration.AutoMigrateAll(db)

	seed.EnsurePermissions(db)
	seed.EnsureBuiltinRoles(db)

	internalseed.RequiredUsers(db)

	return db
}

func constructConnectionString(dbCfg *util.DatabaseConfig) string {
	return "host=" + dbCfg.Host + " user=" + dbCfg.User + " password=" + dbCfg.Password + " dbname=" + dbCfg.Database + " port=" + dbCfg.Port
}
