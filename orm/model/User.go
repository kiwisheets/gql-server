package model

import (
	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/permission"
)

// `User` belongs to `Company`, `CompanyID` is the foreign key

// User model
type User struct {
	SoftDelete
	Email            string `gorm:"UNIQUE_INDEX:idx_email"`
	Phone            *string
	Mobile           *string
	PreferredContact *PreferredContact
	Password         string
	Firstname        string
	Lastname         string
	CompanyID        hide.ID                  `json:"company"`
	Company          Company                  `json:"-"`
	BuiltinRoles     []permission.BuiltinRole `gorm:"many2many:user_builtinroles" json:"-"`
	CustomRoles      []permission.CustomRole  `gorm:"many2many:user_customroles" json:"-"`
	TwoFactor        TwoFactor
}

func (User) IsEntity() {}

// Roles returns all roles a user has as the Role interface
func (u User) Roles() []permission.Role {
	roles := make([]permission.Role, len(u.BuiltinRoles)+len(u.CustomRoles))

	var index uint

	for _, role := range u.BuiltinRoles {
		roles[index] = role
		index++
	}

	for _, role := range u.CustomRoles {
		roles[index] = role
		index++
	}

	return roles
}
