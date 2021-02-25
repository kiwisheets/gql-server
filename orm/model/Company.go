package model

// `Company` has many `Users`

const CompanyBillingAddressType = "company_billing"
const CompanyShippingAddressType = "company_shipping"

// Company model
type Company struct {
	SoftDelete
	Code    string `gorm:"unique_index:idx_code"`
	Name    string
	Users   []User
	Domains []Domain
	Website string
	Clients []Client

	BillingAddress  Address `gorm:"polymorphic:Addressee;polymorphicValue:company_billing"`
	ShippingAddress Address `gorm:"polymorphic:Addressee;polymorphicValue:company_shipping"`
}

func (Company) IsEntity() {}
