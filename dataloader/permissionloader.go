package dataloader

import (
	"time"

	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/permission"
	"github.com/kiwisheets/gql-server/dataloader/generated"
	"gorm.io/gorm"
)

type Result struct {
	permission.Permission
	UserID hide.ID
}

func newPermissionsLoaderByUserIDLoader(db *gorm.DB) *generated.PermissionsLoader {
	return generated.NewPermissionsLoader(generated.PermissionsLoaderConfig{
		MaxBatch: 500,
		Wait:     1 * time.Millisecond,
		Fetch: func(userIDs []int64) ([][]*permission.Permission, []error) {
			permissionsByUserID := make(map[int64][]*permission.Permission)

			rows, err := db.Raw(`
				SELECT DISTINCT user_builtinroles.user_id AS user_id, "id", subject, operation, "name", description
				FROM user_builtinroles
				LEFT JOIN builtinrole_permissions ON builtinrole_permissions.builtin_role_id = user_builtinroles.builtin_role_id
				LEFT JOIN user_customroles ON user_customroles.user_id = user_builtinroles.user_id
				LEFT JOIN customrole_permissions ON customrole_permissions.custom_role_id = user_customroles.custom_role_id
				LEFT JOIN permissions ON permissions.id = builtinrole_permissions.permission_id OR permissions.id = customrole_permissions.permission_id
				WHERE user_builtinroles.user_id IN (?)
			`, userIDs).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			for rows.Next() {
				var result Result
				db.ScanRows(rows, &result)
				if result.ID == 0 {
					// no value returned
				} else {
					permissionsByUserID[int64(result.UserID)] =
						append(permissionsByUserID[int64(result.UserID)], &result.Permission)
				}
			}

			orderedPermissions := make([][]*permission.Permission, len(userIDs))
			var errs []error
			for i, userID := range userIDs {
				orderedPermissions[i] = permissionsByUserID[userID]
			}

			return orderedPermissions, errs
		},
	})
}

// func newPermissionCheckerByIDLoader(db *gorm.DB) *generated.PermissionChecker {
// 	return generated.NewPermissionChecker(generated.PermissionCheckerConfig{
// 		MaxBatch: 200,
// 		Wait: 1 * time.Millisecond,
// 		Fetch: func(keys []int64) ([]*bool, []error) {
// 			success := make([]*bool, len(keys))
// 			errors := make([]error, len(keys))

// 			for i, key := range keys {
// 				var builtinRoles []*model.BuiltinRole
// 				err := db.Model(&model.User{
// 					ModelSoftDelete: model.ModelSoftDelete{
// 						ID: hide.ID(key),
// 					},
// 				}).Related(builtinRoles, "BuiltinRoles")

// 				if err != nil {
// 					errors[i] = err
// 					continue
// 				}

// 				for _, role := range builtinRoles {

// 				}
// 			}
// 		}
// 	})
// }
