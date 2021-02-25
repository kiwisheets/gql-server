package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/emvi/hide"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/orm/model"
)

func (r *entityResolver) FindClientByID(ctx context.Context, id hide.ID) (*model.Client, error) {
	client := &model.Client{}

	err := r.DB.Where(id).Find(&client).Error

	return client, err
}

func (r *entityResolver) FindCompanyByID(ctx context.Context, id hide.ID) (*model.Company, error) {
	company := &model.Company{}

	err := r.DB.Where(id).Find(&company).Error

	return company, err
}

func (r *entityResolver) FindUserByID(ctx context.Context, id hide.ID) (*model.User, error) {
	user, err := dataloader.For(ctx).UserByID.Load(int64(id))

	return user, err
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
