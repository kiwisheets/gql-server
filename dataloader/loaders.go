// Package dataloader contains efficient dataloaders
package dataloader

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
)

var loadersKey = &contextKey{"dataloaders"}

type contextKey struct {
	name string
}

// Loaders structure contains usable dataloaders
type Loaders struct {
	UserByID                        *generated.UserLoader
	UsersByCompanyID                *generated.UserSliceLoader
	UserByEmail                     *generated.UserStringLoader
	CompanyByID                     *generated.CompanyLoader
	CompanyByUserID                 *generated.CompanyLoader
	CompanyByCode                   *generated.CompanyStringLoader
	DomainsByCompanyID              *generated.DomainSliceLoader
	RolesByUserID                   *generated.RoleLoader
	PermissionsByUserID             *generated.PermissionsLoader
	ClientBillingAddressByClientID  *generated.AddressLoader
	ClientShippingAddressByClientID *generated.AddressLoader
	//PermissionByUserID *generated.PermissionsLoader
}

// Middleware handles dataloader requests
func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		loaders := &Loaders{
			UserByID:                        newUserByIDLoader(db, ctx),
			UsersByCompanyID:                newUsersByCompanyIDLoader(db),
			UserByEmail:                     newUserByEmailLoader(db),
			CompanyByID:                     newCompanyByIDLoader(db),
			CompanyByUserID:                 newCompanyByUserIDLoader(db),
			CompanyByCode:                   newCompanyByCodeLoader(db),
			DomainsByCompanyID:              newDomainsByCompanyIDLoader(db),
			RolesByUserID:                   newRoleByUserIDLoader(db),
			PermissionsByUserID:             newPermissionsLoaderByUserIDLoader(db),
			ClientBillingAddressByClientID:  newAddressByAddresseeIDLoader(db, model.ClientBillingAddressType),
			ClientShippingAddressByClientID: newAddressByAddresseeIDLoader(db, model.ClientShippingAddressType),
			//PermissionByUserID: newPermissionsByUserIDLoader(db),
		}

		ctx = context.WithValue(
			ctx,
			loadersKey,
			loaders,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// For returns the available dataloaders
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
