package model

import (
	"log"
	"strings"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/operation"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model/permission/subject"
)

// Permission model
type Permission struct {
	Model
	Name        string
	Description string
	Subject     subject.Subject     `gorm:"type:integer"`
	Operation   operation.Operation `gorm:"type:integer"`
}

// CheckPermission check a full permission string (like ones from the graphql schema)
// And return whether or not the permission matches
func (p Permission) CheckPermission(sub subject.Subject, op operation.Operation) bool {
	return (p.Subject == subject.SubjectAny || p.Subject == sub) &&
		(p.Operation == operation.Any || p.Operation == op)
}

func (p Permission) CheckPermissionString(requestedPerm string) bool {
	strings := strings.Split(requestedPerm, ":")
	if len(strings) != 2 {
		return false
	}

	var op operation.Operation

	sub, err := subject.ParseSubject(strings[0])
	if err != nil {
		log.Printf("failed to check permission. unable to parse subject: %s \n", strings[0])
		return false
	}

	op.FromString(strings[1])

	return p.CheckPermission(sub, op)
}

// func (p Permission) CheckPermission(permString string) bool {
// 	strings := strings.Split(permString, ":")
// 	if len(strings) != 2 {
// 		return false
// 	}

// 	var sub subject.Subject
// 	var op operation.Operation

// 	sub.Scan(strings[0])
// 	op.Scan(strings[1])

// 	return (p.Subject == subject.Any || p.Subject == sub) &&
// 		(p.Operation == operation.Any || p.Operation == op)
// }
