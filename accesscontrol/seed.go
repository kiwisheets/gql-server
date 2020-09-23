package accesscontrol

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/operation"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/subject"
	"gorm.io/gorm"
)

func seedPermissions(db *gorm.DB) {
	// All Permissions
	createOrGetPerm(db, subject.SubjectAny, operation.Any)

	// User object Permissions

	// Me permission, allows actions related to the logged in user
	createOrGetPerm(db, subject.SubjectMe, operation.Read)
	createOrGetPerm(db, subject.SubjectMe, operation.Update)

	// User permission, allows actions on a single user
	createOrGetPerm(db, subject.SubjectUser, operation.Create)
	createOrGetPerm(db, subject.SubjectUser, operation.Read)
	createOrGetPerm(db, subject.SubjectUser, operation.Update)
	createOrGetPerm(db, subject.SubjectUser, operation.Delete)

	// Users permission, allows actions on groups of users
	createOrGetPerm(db, subject.SubjectUsers, operation.Read)

	// Company permissions
	createOrGetPerm(db, subject.SubjectCompany, operation.Create)
	createOrGetPerm(db, subject.SubjectCompany, operation.Read)
	createOrGetPerm(db, subject.SubjectCompany, operation.Update)
	createOrGetPerm(db, subject.SubjectCompany, operation.Delete)

	// Client permissions
	createOrGetPerm(db, subject.SubjectClient, operation.Create)
	createOrGetPerm(db, subject.SubjectClient, operation.Read)
	createOrGetPerm(db, subject.SubjectClient, operation.Update)
	createOrGetPerm(db, subject.SubjectClient, operation.Delete)

	// Clients permissions
	createOrGetPerm(db, subject.SubjectClients, operation.Read)

	// Contact permissions
	createOrGetPerm(db, subject.SubjectContact, operation.Create)
	createOrGetPerm(db, subject.SubjectContact, operation.Read)
	createOrGetPerm(db, subject.SubjectContact, operation.Update)
	createOrGetPerm(db, subject.SubjectContact, operation.Delete)

	// Contacts permissions
	createOrGetPerm(db, subject.SubjectContacts, operation.Read)

	// roles

	// service admin role
	db.Where(model.BuiltinRole{
		Name: "Service Admin",
	}).Attrs(model.BuiltinRole{
		Description: "",
		Permissions: []model.Permission{
			createOrGetPerm(db, subject.SubjectAny, operation.Any),
		},
	}).FirstOrCreate(&model.BuiltinRole{})

	// standard user role
	db.Where(model.BuiltinRole{
		Name: "Standard User",
	}).Attrs(model.BuiltinRole{
		Description: "",
		Permissions: []model.Permission{
			createOrGetPerm(db, subject.SubjectMe, operation.Any),
			createOrGetPerm(db, subject.SubjectUser, operation.Read),
			createOrGetPerm(db, subject.SubjectUsers, operation.Read),
			createOrGetPerm(db, subject.SubjectCompany, operation.Read),
			createOrGetPerm(db, subject.SubjectClient, operation.Read),
			createOrGetPerm(db, subject.SubjectClients, operation.Read),
			createOrGetPerm(db, subject.SubjectUserContact, operation.Read),
		},
	}).FirstOrCreate(&model.BuiltinRole{})

	// standard user role
	db.Where(model.CustomRole{
		Name: "Custom User",
	}).Attrs(model.CustomRole{
		Description: "",
		Permissions: []model.Permission{
			createOrGetPerm(db, subject.SubjectContacts, operation.Read),
			createOrGetPerm(db, subject.SubjectClientContact, operation.Read),
			createOrGetPerm(db, subject.SubjectClient, operation.Create),
		},
		CompanyID: 1,
	}).FirstOrCreate(&model.CustomRole{})

}
