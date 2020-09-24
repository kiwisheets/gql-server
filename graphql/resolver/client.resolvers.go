package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/modelgen"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/emvi/hide"
)

func (r *clientResolver) BillingAddress(ctx context.Context, obj *model.Client) (*model.Address, error) {
	return dataloader.For(ctx).ClientBillingAddressByClientID.Load(obj.IDint())
}

func (r *clientResolver) ShippingAddress(ctx context.Context, obj *model.Client) (*model.Address, error) {
	return dataloader.For(ctx).ClientShippingAddressByClientID.Load(obj.IDint())
}

func (r *clientResolver) Contacts(ctx context.Context, obj *model.Client) ([]*model.Contact, error) {
	preferredContact := model.PreferredContactEmail

	return []*model.Contact{{
		SoftDelete: model.SoftDelete{
			ID:        1,
			CreatedAt: time.Now(),
		},
		Email:            util.String("email@email.com"),
		Firstname:        "first name",
		Lastname:         "last name",
		Mobile:           util.String("0123456789"),
		PreferredContact: &preferredContact,
	}}, nil
}

func (r *mutationResolver) CreateClient(ctx context.Context, client modelgen.CreateClientInput) (*model.Client, error) {
	clientObject := model.Client{
		Name:           client.Name,
		CompanyID:      auth.For(ctx).User.CompanyID,
		Phone:          client.Phone,
		Website:        client.Website,
		VatNumber:      client.VatNumber,
		BusinessNumber: client.BusinessNumber,
	}

	if client.BillingAddress != nil {
		clientObject.BillingAddress = &model.Address{
			Name:       client.BillingAddress.Name,
			Street1:    client.BillingAddress.Street1,
			Street2:    client.BillingAddress.Street2,
			City:       client.BillingAddress.City,
			PostalCode: client.BillingAddress.PostalCode,
			State:      client.BillingAddress.State,
			Country:    client.BillingAddress.Country,
		}
	}

	if client.ShippingAddress != nil {
		clientObject.ShippingAddress = &model.Address{
			Name:       client.ShippingAddress.Name,
			Street1:    client.ShippingAddress.Street1,
			Street2:    client.ShippingAddress.Street2,
			City:       client.ShippingAddress.City,
			PostalCode: client.ShippingAddress.PostalCode,
			State:      client.ShippingAddress.State,
			Country:    client.ShippingAddress.Country,
		}
	}

	// changeset.ApplyChanges(, &clientObject)

	err := r.DB.Create(&clientObject).Error

	if clientObject.ID == 0 || err != nil {
		return nil, fmt.Errorf("Unable to create Client")
	}

	return &clientObject, nil
}

func (r *mutationResolver) UpdateClient(ctx context.Context, client modelgen.UpdateClientInput) (*model.Client, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteClient(ctx context.Context, id hide.ID) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Client(ctx context.Context, id hide.ID) (*model.Client, error) {
	var client model.Client
	r.DB.Model(&model.Client{
		SoftDelete: model.SoftDelete{
			ID: 1,
		},
	}).First(&client)
	if client.ID == 0 {
		return nil, fmt.Errorf("Client not found")
	}

	return &client, nil
}

func (r *queryResolver) Clients(ctx context.Context, page *int) ([]*model.Client, error) {
	limit := 20
	clients := make([]*model.Client, limit)
	if page == nil {
		page = util.Int(0)
	}
	r.DB.Order("name").Limit(limit).Offset(limit * *page).Find(&clients)

	return clients, nil
}

// Client returns generated.ClientResolver implementation.
func (r *Resolver) Client() generated.ClientResolver { return &clientResolver{r} }

type clientResolver struct{ *Resolver }
