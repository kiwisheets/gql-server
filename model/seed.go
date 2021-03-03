package model

import (
	"github.com/kiwisheets/auth/permission"
	"github.com/kiwisheets/gql-server/auth/password"
	"github.com/kiwisheets/gql-server/util"
	"gorm.io/gorm"
)

// RequiredUsers ensures that at least a default ServiceAdmin account exists
func RequiredUsers(db *gorm.DB) {

	var company Company

	db.Where(Company{
		Code: "sa",
	}).Attrs(Company{
		Name:    "Service Admins",
		Website: "https://kiwisheets.com",
		BillingAddress: Address{
			Name:       "Test",
			Street1:    "123 Some Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
		ShippingAddress: Address{
			Name:       "Test",
			Street1:    "123 Some Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
	}).FirstOrCreate(&company)

	var domain Domain

	// TODO: make default domain configurable via env variables
	db.Where(Domain{
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

	var user User

	db.Where(User{
		CompanyID: company.ID,
		// Check role
	}).Attrs(User{
		Email:     "serviceadmin@" + domain.Domain,
		Firstname: "Service",
		Lastname:  "Admin",
		Password:  hash,
		BuiltinRoles: []permission.BuiltinRole{
			serviceAdminRole,
		},
	}).FirstOrCreate(&user)

	var secondUser User
	hash, _ = password.HashPassword("password")

	db.Where(User{
		CompanyID: company.ID,
		Email:     "testuser@" + domain.Domain,
		// Check role
	}).Attrs(User{
		Firstname: "Test",
		Lastname:  "User",
		Password:  hash,
		BuiltinRoles: []permission.BuiltinRole{
			standardUserRole,
		},
	}).FirstOrCreate(&secondUser)

	var demoClient Client

	db.Where(Client{
		Name:      "3B",
		CompanyID: company.ID,
	}).Attrs(Client{
		Phone:          util.String("+6421456789"),
		VatNumber:      util.String("1234567890"),
		BusinessNumber: util.String("0987654321"),
		Website:        util.String("https://website.com"),
		BillingAddress: Address{
			Name:       "Aaron",
			Street1:    "123 Make Believe Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
		ShippingAddress: Address{
			Name:       "Aaron",
			Street1:    "123 Make Believe Street",
			City:       "Auckland",
			PostalCode: 1234,
			Country:    "New Zealand",
		},
	}).FirstOrCreate(&demoClient)
}
