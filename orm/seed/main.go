package seed

import (
	"github.com/kiwisheets/auth/permission"
	password "github.com/kiwisheets/gql-server/auth"
	"github.com/kiwisheets/gql-server/orm/model"
	"github.com/kiwisheets/gql-server/util"
	"gorm.io/gorm"
)

// RequiredUsers ensures that at least a default ServiceAdmin account exists
func RequiredUsers(db *gorm.DB) {

	var company model.Company

	db.Where(model.Company{
		Code: "sa",
	}).Attrs(model.Company{
		Name:    "Service Admins",
		Website: "https://kiwisheets.com",
		BillingAddress: model.Address{
			Name:       "Test",
			Street1:    "123 Some Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
		ShippingAddress: model.Address{
			Name:       "Test",
			Street1:    "123 Some Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
	}).FirstOrCreate(&company)

	var domain model.Domain

	// TODO: make default domain configurable via env variables
	db.Where(model.Domain{
		Domain:    "kiwisheets.com",
		CompanyID: company.ID,
	}).FirstOrCreate(&domain)

	// TODO: make default password configurable via env variables
	hash, _ := password.HashPassword("servicepass")

	var serviceAdminRole permission.BuiltinRole
	var standardUserRole permission.BuiltinRole

	// get service admin role
	db.Where(permission.BuiltinRole{
		Name: "Service Admin",
	}).First(&serviceAdminRole)

	db.Where(permission.BuiltinRole{
		Name: "Standard User",
	}).First(&standardUserRole)

	var user model.User

	db.Where(model.User{
		CompanyID: company.ID,
		// Check role
	}).Attrs(model.User{
		Email:     "serviceadmin@" + domain.Domain,
		Firstname: "Service",
		Lastname:  "Admin",
		Password:  hash,
		BuiltinRoles: []permission.BuiltinRole{
			serviceAdminRole,
		},
	}).FirstOrCreate(&user)

	var secondUser model.User
	hash, _ = password.HashPassword("password")

	db.Where(model.User{
		CompanyID: company.ID,
		Email:     "testuser@" + domain.Domain,
		// Check role
	}).Attrs(model.User{
		Firstname: "Test",
		Lastname:  "User",
		Password:  hash,
		BuiltinRoles: []permission.BuiltinRole{
			standardUserRole,
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
		BillingAddress: model.Address{
			Name:       "Aaron",
			Street1:    "123 Make Believe Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
		ShippingAddress: model.Address{
			Name:       "Aaron",
			Street1:    "123 Make Believe Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
	}).FirstOrCreate(&demoClient)
}
