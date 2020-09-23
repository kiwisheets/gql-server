package dataloader

import (
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"github.com/emvi/hide"
	"gorm.io/gorm"
)

func newRoleByUserIDLoader(db *gorm.DB) *generated.RoleLoader {
	return generated.NewRoleLoader(generated.RoleLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(userIDs []int64) ([][]model.Role, []error) {
			// db.Model(&model.)

			db.Model(&model.BuiltinRole{})

			roles := make([][]model.Role, len(userIDs))
			errors := make([]error, len(userIDs))

			for i, userID := range userIDs {
				{
					var builtinRoles []model.BuiltinRole
					err := db.Model(&model.User{
						SoftDelete: model.SoftDelete{
							ID: hide.ID(userID),
						},
					}).Preload("Permissions").Association("BuiltinRoles").Find(&builtinRoles)

					// db.Table("user_builtinroles").Select("builtin_role_id").Where("user_id", userID)

					if err != nil {
						errors[i] = err
						continue
					}

					for _, role := range builtinRoles {
						roles[i] = append(roles[i], role)
					}
				}
				var customRoles []*model.CustomRole
				{
					err := db.Model(&model.User{
						SoftDelete: model.SoftDelete{
							ID: hide.ID(userID),
						},
					}).Preload("Permissions").Association("CustomRoles").Find(&customRoles)

					if err != nil {
						errors[i] = err
						continue
					}
					for _, role := range customRoles {
						roles[i] = append(roles[i], role)
					}
				}
			}

			return roles, errors
		},
	})
}
