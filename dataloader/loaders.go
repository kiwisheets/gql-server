// Package dataloader contains efficient dataloaders
package dataloader

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kiwisheets/gql-server/dataloader/generated"
	"github.com/kiwisheets/gql-server/model"
)

var loadersKey = &dataloaderContextKey{"dataloaders"}

type dataloaderContextKey struct {
	name string
}

// Loaders structure contains usable dataloaders
type Loaders struct {
	UserByID                          *generated.UserLoader
	UsersByCompanyID                  *generated.UserSliceLoader
	UserByEmail                       *generated.UserStringLoader
	CompanyByID                       *generated.CompanyLoader
	CompanyByUserID                   *generated.CompanyLoader
	CompanyByCode                     *generated.CompanyStringLoader
	DomainsByCompanyID                *generated.DomainSliceLoader
	RolesByUserID                     *generated.RoleLoader
	PermissionsByUserID               *generated.PermissionsLoader
	ClientBillingAddressByClientID    *generated.AddressLoader
	ClientShippingAddressByClientID   *generated.AddressLoader
	CompanyBillingAddressByCompanyID  *generated.AddressLoader
	CompanyShippingAddressByCompanyID *generated.AddressLoader
	//PermissionByUserID *generated.PermissionsLoader
}

// Middleware handles dataloader requests
func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		loaders := &Loaders{
			UserByID:                          newUserByIDLoader(db, ctx),
			UsersByCompanyID:                  newUsersByCompanyIDLoader(db),
			UserByEmail:                       newUserByEmailLoader(db),
			CompanyByID:                       newCompanyByIDLoader(db),
			CompanyByUserID:                   newCompanyByUserIDLoader(db),
			CompanyByCode:                     newCompanyByCodeLoader(db),
			DomainsByCompanyID:                newDomainsByCompanyIDLoader(db),
			RolesByUserID:                     newRoleByUserIDLoader(db),
			PermissionsByUserID:               newPermissionsByUserIDLoader(db),
			ClientBillingAddressByClientID:    newAddressByAddresseeIDLoader(db, model.ClientBillingAddressType),
			ClientShippingAddressByClientID:   newAddressByAddresseeIDLoader(db, model.ClientShippingAddressType),
			CompanyBillingAddressByCompanyID:  newAddressByAddresseeIDLoader(db, model.CompanyBillingAddressType),
			CompanyShippingAddressByCompanyID: newAddressByAddresseeIDLoader(db, model.CompanyShippingAddressType),
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
