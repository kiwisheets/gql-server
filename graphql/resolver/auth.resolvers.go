package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/auth/permission"
	internalauth "github.com/kiwisheets/gql-server/auth"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/orm/model"
)

func (r *mutationResolver) Login(ctx context.Context, email string, password string, twoFactor *string) (*model.AuthData, error) {
	// TODO: move logic into auth package

	user, err := dataloader.For(ctx).UserByEmail.Load(email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("Email or Password Incorrect")
	}

	permissions, err := dataloader.For(ctx).PermissionsByUserID.Load(user.IDint())
	if err != nil {
		return nil, err
	}

	if !internalauth.VerifyPassword(r.DB, user.ID, password) {
		return nil, fmt.Errorf("Email or Password Incorrect")
	}

	twoFactorObject, err := internalauth.GetTwoFactor(r.DB, user.ID)

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
		if !internalauth.VerifyTwoFactor(twoFactorObject, *twoFactor) {
			return nil, fmt.Errorf("Invalid 2FA code")
		}
	}

	token, err := internalauth.LoginUser(user, permissions, &r.Cfg.JWT)

	if err != nil {
		return nil, fmt.Errorf("Email or Password Incorrect")
	}

	return &model.AuthData{
		User:             user,
		Token:            &token,
		TwoFactorEnabled: twoFactorObject != nil,
	}, nil
}

func (r *mutationResolver) LoginSecure(ctx context.Context, password string) (string, error) {
	user, err := dataloader.For(ctx).UserByID.Load(int64(auth.For(ctx).UserID))
	if err != nil || user == nil {
		return "", fmt.Errorf("Email or Password Incorrect")
	}
	permissions, err := dataloader.For(ctx).PermissionsByUserID.Load(user.IDint())
	if err != nil {
		return "", err
	}

	return internalauth.LoginUserSecure(user, permissions, &r.Cfg.JWT)
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	userID := auth.For(ctx).UserID

	if !internalauth.VerifyPassword(r.DB, userID, oldPassword) {
		return false, fmt.Errorf("Password Incorrect")
	}

	hash, err := internalauth.HashPassword(newPassword)
	if err != nil {
		return false, fmt.Errorf("Bad password")
	}

	if err := r.DB.Model(&model.User{}).Where(userID).Update("Password", hash).Error; err != nil {
		return true, err
	}

	return true, nil
}

func (r *mutationResolver) NewTwoFactorBackups(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EnableTwoFactor(ctx context.Context, secret string, token string) ([]string, error) {
	return internalauth.EnableTwoFactor(r.DB, auth.For(ctx).UserID, secret, token)
}

func (r *mutationResolver) DisableTwoFactor(ctx context.Context, password string) (bool, error) {
	return internalauth.DisableTwoFactor(r.DB, auth.For(ctx).UserID, password)
}

func (r *permissionResolver) Subject(ctx context.Context, obj *permission.Permission) (string, error) {
	return obj.Subject.String(), nil
}

func (r *permissionResolver) Operation(ctx context.Context, obj *permission.Permission) (string, error) {
	return obj.Operation.String(), nil
}

func (r *queryResolver) TwoFactorBackups(ctx context.Context) ([]string, error) {
	return internalauth.GetBackupKeys(r.DB, auth.For(ctx))
}

func (r *queryResolver) TwoFactorEnabled(ctx context.Context) (bool, error) {
	return internalauth.IsTwoFactorEnabled(r.DB, auth.For(ctx).UserID)
}

func (r *queryResolver) Scopes(ctx context.Context) ([]*permission.Permission, error) {
	permissions, err := dataloader.For(ctx).PermissionsByUserID.Load(int64(auth.For(ctx).UserID))
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Permission returns generated.PermissionResolver implementation.
func (r *Resolver) Permission() generated.PermissionResolver { return &permissionResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type permissionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
