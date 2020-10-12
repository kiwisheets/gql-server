package seed

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"gorm.io/gorm"
)

// RequiredUsers ensures that at least a default ServiceAdmin account exists
func RequiredUsers(db *gorm.DB) {
	var company model.Company

	db.Where(model.Company{
		Code: "sa",
	}).Attrs(model.Company{
		Name: "Service Admins",
	}).FirstOrCreate(&company)

	var domain model.Domain

	// TODO: make default domain configurable via env variables
	db.Where(model.Domain{
		Domain:    "kiwisheets.com",
		CompanyID: company.ID,
	}).FirstOrCreate(&domain)

	// TODO: make default password configurable via env variables
	hash, _ := auth.HashPassword("servicepass")

	var serviceAdminRole model.BuiltinRole
	var standardUserRole model.BuiltinRole
	var customUserRole model.CustomRole

	// get service admin role
	db.Where(model.BuiltinRole{
		Name: "Service Admin",
	}).First(&serviceAdminRole)

	db.Where(model.BuiltinRole{
		Name: "Standard User",
	}).First(&standardUserRole)

	db.Where(model.CustomRole{
		Name: "Custom User",
	}).First(&customUserRole)

	var user model.User

	db.Where(model.User{
		CompanyID: company.ID,
		// Check role
	}).Attrs(model.User{
		Email:     "serviceadmin@" + domain.Domain,
		Firstname: "Service",
		Lastname:  "Admin",
		Password:  hash,
		BuiltinRoles: []model.BuiltinRole{
			serviceAdminRole,
		},
	}).FirstOrCreate(&user)

	var secondUser model.User
	hash, _ = auth.HashPassword("password")

	db.Where(model.User{
		CompanyID: company.ID,
		Email:     "testuser@" + domain.Domain,
		// Check role
	}).Attrs(model.User{
		Firstname: "Test",
		Lastname:  "User",
		Password:  hash,
		BuiltinRoles: []model.BuiltinRole{
			standardUserRole,
		},
		CustomRoles: []model.CustomRole{
			customUserRole,
		},
	}).FirstOrCreate(&secondUser)

	var demoClient model.Client

	db.Where(model.Client{
		Name:      "3B",
		CompanyID: company.ID,
	}).Attrs(model.Client{
		Phone:          util.String("+6421456789"),
		VatNumber:      util.String("1234567890"),
		BusinessNumber: util.String("0987654321"),
		Website:        util.String("https://website.com"),
	}).FirstOrCreate(&demoClient)
}
