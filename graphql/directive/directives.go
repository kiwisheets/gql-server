/*
Package directive implements directives for the graphql schema
*/
package directive

import (
	"context"
	"fmt"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/99designs/gqlgen/graphql"
	"gorm.io/gorm"
)

// Register function registers all directives
func Register(db *gorm.DB, cfg *util.ServerConfig) generated.DirectiveRoot {
	return generated.DirectiveRoot{
		IsAuthenticated: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
			if auth.For(ctx).User == nil {
				return nil, fmt.Errorf("not logged in")
			}
			return next(ctx)
		},
		IsSecureAuthenticated: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
			if auth.For(ctx).User == nil {
				return nil, fmt.Errorf("not logged in")
			}
			if !auth.For(ctx).Secure {
				return nil, fmt.Errorf("not logged in with time sensitive token")
			}
			return next(ctx)
		},
		HasPerm: func(ctx context.Context, obj interface{}, next graphql.Resolver, perm string) (res interface{}, err error) {
			if auth.For(ctx).User == nil {
				return nil, fmt.Errorf("not logged in")
			}

			// should probably optimise this directive as it will be called on most requests
			// roles, err := dataloader.For(ctx).RolesByUserID.Load(auth.For(ctx).User.IDint())

			permissions, err := dataloader.For(ctx).PermissionsByUserID.Load(auth.For(ctx).User.IDint())

			if err != nil {
				return nil, err
			}

			for _, p := range permissions {
				if p.CheckPermissionString(perm) {
					return next(ctx)
				}
			}

			// for _, r := range roles {
			// 	if r.CheckPermission(perm) {
			// 		// permission passed
			// 		return next(ctx)
			// 	}
			// }

			return nil, fmt.Errorf("not authorised")
		},
		HasPerms: func(ctx context.Context, obj interface{}, next graphql.Resolver, requestedPerms []string) (res interface{}, err error) {
			if auth.For(ctx).User == nil {
				return nil, fmt.Errorf("not logged in")
			}

			permsPassed := make([]bool, len(requestedPerms))

			permissions, err := dataloader.For(ctx).PermissionsByUserID.Load(auth.For(ctx).User.IDint())

			if err != nil {
				return nil, err
			}

			for _, userPerm := range permissions {
				for i, requestedPerm := range requestedPerms {
					if permsPassed[i] {
						continue
					}
					if userPerm.CheckPermissionString(requestedPerm) {
						permsPassed[i] = true
					}
				}
			}

			// // should probably optimise this directive as it will be called on most requests
			// roles, err := dataloader.For(ctx).RolesByUserID.Load(auth.For(ctx).User.IDint())

			// if err != nil {
			// 	return nil, err
			// }

			// for _, r := range roles {
			// 	for i, p := range perms {
			// 		if permsPassed[i] {
			// 			continue
			// 		}
			// 		if r.CheckPermission(p) {
			// 			// permission passed
			// 			permsPassed[i] = true
			// 		}
			// 	}
			// }

			for _, p := range permsPassed {
				if !p {
					return nil, fmt.Errorf("not authorised")
				}
			}

			return next(ctx)
		},
	}
}
