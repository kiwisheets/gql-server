package accesscontrol

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/operation"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/subject"
	"github.com/jinzhu/gorm"
)

func seedPermissions(db *gorm.DB) {
	// All Permissions
	newPerm(db, subject.Any, operation.Any)

	// User object Permissions

	// Me permission, allows actions related to the logged in user
	newPerm(db, subject.Me, operation.Read)
	newPerm(db, subject.Me, operation.Update)

	// User permission, allows actions on a single user
	newPerm(db, subject.User, operation.Create)
	newPerm(db, subject.User, operation.Read)
	newPerm(db, subject.User, operation.Update)
	newPerm(db, subject.User, operation.Delete)

	// Users permission, allows actions on groups of users
	newPerm(db, subject.Users, operation.Read)

	// Company permissions
	newPerm(db, subject.Company, operation.Create)
	newPerm(db, subject.Company, operation.Read)
	newPerm(db, subject.Company, operation.Update)
	newPerm(db, subject.Company, operation.Delete)
}
