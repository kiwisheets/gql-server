package migration

import (
	"log"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"github.com/kiwisheets/auth/permission"
	"gorm.io/gorm"
)

func AutoMigrateAll(db *gorm.DB) {

	db.DisableForeignKeyConstraintWhenMigrating = true
	log.Println("Migrating models...")
	db.AutoMigrate(&model.Domain{})
	db.AutoMigrate(&model.Company{})

	db.AutoMigrate(&permission.BuiltinRole{})
	db.AutoMigrate(&permission.CustomRole{})
	db.AutoMigrate(&permission.Permission{})

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

	db.Migrator().DropTable(&permission.BuiltinRole{})
	db.Migrator().DropTable(&permission.CustomRole{})
	db.Migrator().DropTable(&permission.Permission{})

	db.Migrator().DropTable(&model.TwoFactor{})
	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.Client{})
	db.Migrator().DropTable(&model.Contact{})
	db.Migrator().DropTable(&model.Address{})
	log.Println("Done")
	db.DisableForeignKeyConstraintWhenMigrating = false
}
