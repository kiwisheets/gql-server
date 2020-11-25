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
)

func (r *mutationResolver) Login(ctx context.Context, email string, password string, twoFactor *string) (*model.AuthData, error) {
	// TODO: move logic into auth package

	user, err := dataloader.For(ctx).UserByEmail.Load(email)

	if err != nil || user == nil {
		return nil, fmt.Errorf("Email or Password Incorrect")
	}

	if !auth.VerifyPassword(user, password) {
		return nil, fmt.Errorf("Email or Password Incorrect")
	}

	twoFactorObject, err := auth.GetTwoFactor(r.DB, user)

	if twoFactorObject == nil && err == nil {
		// twofactor disabled

	} else {
		// twofactor enabled

		// check if twofactor code is empty
		if twoFactor == nil || *twoFactor == "" {
			return &model.AuthData{
				TwoFactorEnabled: true,
			}, nil
		}

		// verify twofactor code, return if invalid
		if !auth.VerifyTwoFactor(twoFactorObject, *twoFactor) {
			return nil, fmt.Errorf("Invalid 2FA code")
		}
	}

	token, err := auth.LoginUser(user, &r.Cfg.JWT)

	if err != nil {
		return nil, fmt.Errorf("Email of Password Incorrect")
	}

	return &model.AuthData{
		User:             user,
		Token:            &token,
		TwoFactorEnabled: twoFactorObject != nil,
	}, nil
}

func (r *mutationResolver) LoginSecure(ctx context.Context, password string) (string, error) {
	return auth.LoginUserSecure(auth.For(ctx).User, &r.Cfg.JWT)
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	user := auth.For(ctx).User

	if !auth.VerifyPassword(user, oldPassword) {
		return false, fmt.Errorf("Password Incorrect")
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return false, fmt.Errorf("Bad password")
	}

	if err := r.DB.Model(&user).Update("Password", hash).Error; err != nil {
		return true, err
	}

	return true, nil
}

func (r *mutationResolver) NewTwoFactorBackups(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EnableTwoFactor(ctx context.Context, secret string, token string) ([]string, error) {
	return auth.EnableTwoFactor(r.DB, auth.For(ctx).User, secret, token)
}

func (r *mutationResolver) DisableTwoFactor(ctx context.Context, password string) (bool, error) {
	return auth.DisableTwoFactor(r.DB, auth.For(ctx).User, password)
}

func (r *queryResolver) TwoFactorBackups(ctx context.Context) ([]string, error) {
	return auth.GetBackupKeys(r.DB, auth.For(ctx))
}

func (r *queryResolver) TwoFactorEnabled(ctx context.Context) (bool, error) {
	return auth.IsTwoFactorEnabled(r.DB, auth.For(ctx).User)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
