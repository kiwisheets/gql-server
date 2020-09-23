package model

import (
	"github.com/emvi/hide"
)

type Client struct {
	SoftDelete
	Name string

	Website        *string
	VatNumber      *string
	BusinessNumber *string
	Phone          *string

	ShippingAddress *Address `gorm:"polymorphic:Addressee;polymorphicValue:shipping"`
	BillingAddress  *Address `gorm:"polymorphic:Addressee;polymorphicValue:billing"`

	Contacts  []Contact
	CompanyID hide.ID `json:"company"`
	Company   Company `json:"-"`
}
