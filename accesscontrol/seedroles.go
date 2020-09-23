package accesscontrol

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/operation"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/subject"
	"gorm.io/gorm"
)

func seedRoles(db *gorm.DB) {

	// Service Admin role
	// Should have all permissions for now

	perm := model.Permission{}

	db.Where(model.Permission{
		Subject:   subject.SubjectAny,
		Operation: operation.Any,
	}).First(&perm)

	serviceAdminRole := model.BuiltinRole{}

	db.Where(model.BuiltinRole{
		Name: "Service Admin",
	}).Attrs(model.BuiltinRole{
		Description: "",
		Permissions: []model.Permission{
			perm,
		},
	}).FirstOrCreate(&serviceAdminRole)
}
