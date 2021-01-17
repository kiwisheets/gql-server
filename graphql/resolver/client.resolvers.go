package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/emvi/hide"
	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/graphql/modelgen"
	"github.com/kiwisheets/gql-server/orm/model"
	"github.com/kiwisheets/gql-server/util"
	"github.com/kiwisheets/gql-server/util/dereference"
	log "github.com/sirupsen/logrus"
)

func (r *clientResolver) ShippingAddress(ctx context.Context, obj *model.Client) (*model.Address, error) {
	return dataloader.For(ctx).ClientShippingAddressByClientID.Load(obj.IDint())
}

func (r *clientResolver) BillingAddress(ctx context.Context, obj *model.Client) (*model.Address, error) {
	return dataloader.For(ctx).ClientBillingAddressByClientID.Load(obj.IDint())
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
		CompanyID:      auth.For(ctx).CompanyID,
		Phone:          client.Phone,
		Website:        client.Website,
		VatNumber:      client.VatNumber,
		BusinessNumber: client.BusinessNumber,
	}

	clientObject.BillingAddress = model.MapInputToAddress(*client.BillingAddress)
	clientObject.ShippingAddress = model.MapInputToAddress(*client.ShippingAddress)

	// changeset.ApplyChanges(, &clientObject)

	err := r.DB.Create(&clientObject).Error

	if clientObject.ID == 0 || err != nil {
		return nil, fmt.Errorf("Unable to create Client")
	}

	return &clientObject, nil
}

func (r *mutationResolver) UpdateClient(ctx context.Context, id hide.ID, client modelgen.UpdateClientInput) (*model.Client, error) {
	res := r.DB.Model(&model.Client{
		SoftDelete: model.SoftDelete{
			ID: id,
		},
	}).Updates(model.Client{
		Name:           dereference.String(client.Name, ""),
		Phone:          client.Phone,
		VatNumber:      client.VatNumber,
		BusinessNumber: client.BusinessNumber,
		Website:        client.Website,
	})

	if res.RowsAffected == 1 {
		var client model.Client
		r.DB.Model(&model.Client{}).Where(id).First(&client)
		if client.ID == 0 {
			return nil, fmt.Errorf("Client not found")
		}
		return &client, nil
	}

	log.Printf("failed to update client: %d", int64(id))

	return nil, fmt.Errorf("failed to update client")
}

func (r *mutationResolver) DeleteClient(ctx context.Context, id hide.ID) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Client(ctx context.Context, id hide.ID) (*model.Client, error) {
	var client model.Client
	r.DB.Model(&model.Client{}).Where(id).First(&client)
	if client.ID == 0 {
		return nil, fmt.Errorf("Client not found")
	}

	return &client, nil
}

func (r *queryResolver) ClientCount(ctx context.Context) (int, error) {
	var count int64
	r.DB.Model(&model.Client{
		CompanyID: auth.For(ctx).CompanyID,
	}).Count(&count)

	return int(count), nil
}

func (r *queryResolver) Clients(ctx context.Context, page *int) ([]*model.Client, error) {
	limit := 20
	clients := make([]*model.Client, limit)
	if page == nil {
		page = util.Int(0)
	}
	r.DB.Order("name").Limit(limit).Offset(limit * int(*page)).Find(&clients)

	return clients, nil
}

// Client returns generated.ClientResolver implementation.
func (r *Resolver) Client() generated.ClientResolver { return &clientResolver{r} }

type clientResolver struct{ *Resolver }
