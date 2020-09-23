package orm

import (
	"log"
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/accesscontrol"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/migration"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/seed"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init connects to and initialises the database
func Init(cfg *util.ServerConfig) *gorm.DB {
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

	accesscontrol.EnsurePermissions(db)
	accesscontrol.EnsureBuiltinRoles(db)

	seed.RequiredUsers(db)

	return db
}

func constructConnectionString(dbCfg *util.DatabaseConfig) string {
	return "host=" + dbCfg.Host + " user=" + dbCfg.User + " password=" + dbCfg.Password + " dbname=" + dbCfg.Database + " port=" + dbCfg.Port
}
