package dataloader

import (
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"github.com/emvi/hide"
	"gorm.io/gorm"
)

type Result struct {
	model.Permission
	UserID hide.ID
}

func newPermissionsLoaderByUserIDLoader(db *gorm.DB) *generated.PermissionsLoader {
	return generated.NewPermissionsLoader(generated.PermissionsLoaderConfig{
		MaxBatch: 500,
		Wait:     1 * time.Millisecond,
		Fetch: func(userIDs []int64) ([][]*model.Permission, []error) {
			// builtInRows, err := db.Raw(`
			// 	SELECT user_id, permission_id AS id, subject, operation, "name", description
			// 	FROM user_builtinroles
			// 	LEFT JOIN builtinrole_permissions on builtinrole_permissions.builtin_role_id = user_builtinroles.builtin_role_id
			// 	LEFT JOIN permissions on permissions.id = permission_id
			// 	WHERE user_id IN (?)
			// `, userIDs).Rows()

			// if err != nil {
			// 	if builtInRows == nil {
			// 		return nil, []error{err}
			// 	}
			// 	// log error
			// }
			// defer builtInRows.Close()

			// customRows, err := db.Raw(`
			// 	SELECT user_id, permission_id as "id", subject, operation, "name", description
			// 	FROM user_customroles
			// 	LEFT JOIN customrole_permissions on customrole_permissions.custom_role_id = user_customroles.custom_role_id
			// 	LEFT JOIN permissions on permissions.id = permission_id
			// 	WHERE user_id IN (?)
			// `, userIDs).Rows()

			// if err != nil {
			// 	if customRows == nil {
			// 		return nil, []error{err}
			// 	}
			// 	// log error
			// }
			// defer customRows.Close()

			permissionsByUserID := make(map[int64][]*model.Permission)
			// for builtInRows.Next() {
			// 	var result Result
			// 	db.ScanRows(builtInRows, &result)

			// 	if result.ID == 0 {
			// 		// no value returned
			// 	} else {
			// 		permissionsByUserID[int64(result.UserID)] =
			// 			append(permissionsByUserID[int64(result.UserID)], &result.Permission)
			// 	}
			// }
			// for customRows.Next() {
			// 	var result Result
			// 	db.ScanRows(customRows, &result)
			// 	if result.ID == 0 {
			// 		// no value returned
			// 	} else {
			// 		permissionsByUserID[int64(result.UserID)] =
			// 			append(permissionsByUserID[int64(result.UserID)], &result.Permission)
			// 	}
			// }

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

			orderedPermissions := make([][]*model.Permission, len(userIDs))
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
