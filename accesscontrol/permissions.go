package accesscontrol

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/operation"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/subject"
	"gorm.io/gorm"
)

// EnsurePermissions ensures that all permissions exist in the database
func EnsurePermissions(db *gorm.DB) {
	seedPermissions(db)
}

// EnsureBuiltinRoles ensure that all builtin roles exist in the database
func EnsureBuiltinRoles(db *gorm.DB) {
	// seedRoles(db)
}

func createOrGetPerm(db *gorm.DB, subject subject.Subject, operation operation.Operation) model.Permission {
	return createOrGetPermWithDesc(
		db,
		subject,
		operation,
		"Allow "+operation.String()+" operations on "+subject.String()+" resources",
	)
}

func createOrGetPermWithDesc(db *gorm.DB, subject subject.Subject, operation operation.Operation, description string) model.Permission {
	var perm model.Permission
	db.FirstOrCreate(&perm, model.Permission{
		Name:        subject.String() + ":" + operation.String(),
		Description: description,
		Subject:     subject,
		Operation:   operation,
	})

	return perm
}
