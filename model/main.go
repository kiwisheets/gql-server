package model

import (
	"log"

	"github.com/kiwisheets/auth/permission"
	"github.com/kiwisheets/auth/seed"
	"github.com/kiwisheets/orm"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

type Orm struct {
	DB *gorm.DB
}

func (o *Orm) Close() {
	sqlDB, _ := o.DB.DB()
	sqlDB.Close()
}

func Init(cfg *util.DatabaseConfig) *Orm {
	orm := Orm{
		DB: orm.Init(cfg),
	}

	migrate(orm.DB)

	seed.EnsurePermissions(orm.DB)
	seed.EnsureBuiltinRoles(orm.DB)

	RequiredUsers(orm.DB)

	return &orm
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&Company{})

	db.AutoMigrate(&permission.BuiltinRole{})
	db.AutoMigrate(&permission.CustomRole{})
	db.AutoMigrate(&permission.Permission{})

	db.AutoMigrate(&TwoFactor{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Client{})
	db.AutoMigrate(&Contact{})
	db.AutoMigrate(&Address{})
}

func DropAll(db *gorm.DB) {
	db.DisableForeignKeyConstraintWhenMigrating = true
	log.Println("Dropping all tables...")
	db.Migrator().DropTable(&Company{})
	db.Migrator().DropTable(&Domain{})

	db.Migrator().DropTable(&permission.BuiltinRole{})
	db.Migrator().DropTable(&permission.CustomRole{})
	db.Migrator().DropTable(&permission.Permission{})

	db.Migrator().DropTable(&TwoFactor{})
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&Client{})
	db.Migrator().DropTable(&Contact{})
	db.Migrator().DropTable(&Address{})
	log.Println("Done")
	db.DisableForeignKeyConstraintWhenMigrating = false
}
