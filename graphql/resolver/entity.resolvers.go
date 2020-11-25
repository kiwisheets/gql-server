package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"github.com/emvi/hide"
)

func (r *entityResolver) FindClientByID(ctx context.Context, id hide.ID) (*model.Client, error) {
	client := &model.Client{}

	err := r.DB.Where(id).Find(&client).Error

	return client, err
}

func (r *entityResolver) FindUserByID(ctx context.Context, id hide.ID) (*model.User, error) {
	user, err := dataloader.For(ctx).UserByID.Load(int64(id))

	return user, err
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
