/*
Package directive implements directives for the graphql schema
*/
package directive

import (
	"github.com/kiwisheets/auth/directive"
	"github.com/kiwisheets/gql-server/graphql/generated"
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
