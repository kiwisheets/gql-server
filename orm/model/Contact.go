package model

import "github.com/emvi/hide"

type Contact struct {
	SoftDelete
	ClientID         hide.ID
	Email            *string
	Phone            *string
	Mobile           *string
	PreferredContact *PreferredContact
	Firstname        string
	Lastname         string
}
