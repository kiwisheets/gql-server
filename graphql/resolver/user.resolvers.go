package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/emvi/hide"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateUser(ctx context.Context, email string, password string) (*model.User, error) {
	// get code

	company := &auth.For(ctx).User.Company

	// verify that this user has the ability to create a user for this company

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("unable to create User. Password invalid")
	}

	var user = model.User{
		Company:  *company,
		Email:    email,
		Password: hash,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("unable to create User. Already exists")
	}

	return &user, nil
}

func (r *mutationResolver) CreateUserForCompany(ctx context.Context, companyID hide.ID, email string, password string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id hide.ID) (*bool, error) {
	if err := r.DB.Delete(&model.User{
		SoftDelete: model.SoftDelete{
			ID: id,
		},
	}).Error; err == gorm.ErrRecordNotFound {
		return util.Bool(false), fmt.Errorf("User not found")
	} else if err != nil {
		return util.Bool(false), fmt.Errorf("No user specified")
	}

	return util.Bool(true), nil
}

func (r *mutationResolver) DeleteUserForCompany(ctx context.Context, companyID hide.ID, id hide.ID) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, ids []hide.ID) ([]bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUsersForCompany(ctx context.Context, companyID hide.ID, ids []hide.ID) ([]bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return auth.For(ctx).User, nil
}

func (r *queryResolver) User(ctx context.Context, id hide.ID) (*model.User, error) {
	user, err := dataloader.For(ctx).UserByID.Load(int64(id))

	return user, err
}

func (r *queryResolver) UserForCompany(ctx context.Context, companyID hide.ID, id hide.ID) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, page *int) ([]*model.User, error) {
	limit := 20
	users := make([]*model.User, limit)
	if page == nil {
		page = util.Int(0)
	}
	r.DB.Order("firstname").Limit(limit).Offset(limit * *page).Find(&users)

	return users, nil
}

func (r *queryResolver) UsersForCompany(ctx context.Context, companyID hide.ID, page *int) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchUsers(ctx context.Context, search string, page *int) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchUsersForCompany(ctx context.Context, companyID hide.ID, search string, page *int) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Company(ctx context.Context, obj *model.User) (*model.Company, error) {
	return dataloader.For(ctx).CompanyByID.Load(int64(obj.CompanyID))
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
