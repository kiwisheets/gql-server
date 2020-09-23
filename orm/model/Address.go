package model

import "github.com/emvi/hide"

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
