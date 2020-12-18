package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/emvi/hide"
	"github.com/kiwisheets/auth"
	internalauth "github.com/kiwisheets/gql-server/auth"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/orm/model"
	"github.com/kiwisheets/gql-server/util"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateUser(ctx context.Context, email string, password string) (*model.User, error) {
	// get code

	company, err := dataloader.For(ctx).CompanyByID.Load(int64(auth.For(ctx).CompanyID))
	if err != nil {
		return nil, fmt.Errorf("unable to create User. Company not found")
	}

	// verify that this user has the ability to create a user for this company

	hash, err := internalauth.HashPassword(password)
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
	if err := r.DB.Where("company_id = ?", auth.For(ctx).CompanyID).Delete(&model.User{
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
	var user model.User
	if err := r.DB.Model(&model.User{}).Where(auth.For(ctx).UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *queryResolver) User(ctx context.Context, id hide.ID) (*model.User, error) {
	var user model.User
	err := r.DB.Where("company_id = ?", auth.For(ctx).CompanyID).Where(id).First(&user).Error
	return &user, err
}

func (r *queryResolver) UserForCompany(ctx context.Context, companyID hide.ID, id hide.ID) (*model.User, error) {
	var user model.User
	err := r.DB.Where("company_id = ?", companyID).Where(id).First(&user).Error
	return &user, err
}

func (r *queryResolver) Users(ctx context.Context, page *int) ([]*model.User, error) {
	limit := 20
	users := make([]*model.User, limit)
	if page == nil {
		page = util.Int(0)
	}
	err := r.DB.Where("company_id = ?", auth.For(ctx).CompanyID).Order("firstname").Limit(limit).Offset(limit * *page).Find(&users).Error

	return users, err
}

func (r *queryResolver) UsersForCompany(ctx context.Context, companyID hide.ID, page *int) ([]*model.User, error) {
	limit := 20
	users := make([]*model.User, limit)
	if page == nil {
		page = util.Int(0)
	}
	err := r.DB.Where("company_id = ?", companyID).Order("firstname").Limit(limit).Offset(limit * *page).Find(&users).Error

	return users, err
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
