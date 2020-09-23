package migration

import (
	"log"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"gorm.io/gorm"
)

func AutoMigrateAll(db *gorm.DB) {

	db.DisableForeignKeyConstraintWhenMigrating = true
	log.Println("Migrating models...")
	db.AutoMigrate(&model.Domain{})
	db.AutoMigrate(&model.Company{})
	db.AutoMigrate(&model.BuiltinRole{})
	db.AutoMigrate(&model.CustomRole{})
	db.AutoMigrate(&model.Permission{})
	db.AutoMigrate(&model.TwoFactor{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Client{})
	db.AutoMigrate(&model.Contact{})
	db.AutoMigrate(&model.Address{})
	log.Println("Done")
	db.DisableForeignKeyConstraintWhenMigrating = false
}

func DropAll(db *gorm.DB) {
	db.DisableForeignKeyConstraintWhenMigrating = true
	log.Println("Dropping all tables...")
	db.Migrator().DropTable(&model.Company{})
	db.Migrator().DropTable(&model.Domain{})
	db.Migrator().DropTable(&model.BuiltinRole{})
	db.Migrator().DropTable(&model.CustomRole{})
	db.Migrator().DropTable(&model.Permission{})
	db.Migrator().DropTable(&model.TwoFactor{})
	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.Client{})
	db.Migrator().DropTable(&model.Contact{})
	db.Migrator().DropTable(&model.Address{})
	log.Println("Done")
	db.DisableForeignKeyConstraintWhenMigrating = false
}
