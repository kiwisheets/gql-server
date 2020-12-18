/*
Package directive implements directives for the graphql schema
*/
package directive

import (
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"github.com/kiwisheets/auth/directive"
)

// Register function registers all directives
func Register() generated.DirectiveRoot {
	return generated.DirectiveRoot{
		IsAuthenticated:       directive.IsAuthenticated,
		IsSecureAuthenticated: directive.IsSecureAuthenticated,
		HasPerm:               directive.HasPerm,
		HasPerms:              directive.HasPerms,
	}
}
