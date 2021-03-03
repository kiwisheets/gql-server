package model

import (
	"github.com/emvi/hide"
)

const ClientBillingAddressType = "client_billing"
const ClientShippingAddressType = "client_shipping"

type Client struct {
	SoftDelete
	Name string

	Website        *string
	VatNumber      *string
	BusinessNumber *string
	Phone          *string

	BillingAddress  Address `gorm:"polymorphic:Addressee;polymorphicValue:client_billing"`
	ShippingAddress Address `gorm:"polymorphic:Addressee;polymorphicValue:client_shipping"`

	Contacts  []Contact
	CompanyID hide.ID `json:"company"`
	Company   Company `json:"-"`
}

func (Client) IsEntity() {}
