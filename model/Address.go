package model

import (
	"github.com/emvi/hide"
)

type Address struct {
	SoftDelete
	Name       string
	Street1    string
	Street2    *string
	City       string
	State      *string
	PostalCode int
	Country    string

	AddresseeID   hide.ID
	AddresseeType string
}

func MapInputToAddress(in CreateAddressInput) Address {
	return Address{
		Name:       in.Name,
		Street1:    in.Street1,
		Street2:    in.Street2,
		City:       in.City,
		State:      in.State,
		PostalCode: in.PostalCode,
		Country:    in.Country,
	}
}
