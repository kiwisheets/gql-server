package dataloader

import (
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"gorm.io/gorm"
)

func newUserByIDLoader(db *gorm.DB) *generated.UserLoader {
	return generated.NewUserLoader(generated.UserLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []int64) ([]*model.User, []error) {
			rows, err := db.Model(&model.User{}).Where(ids).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			if err != nil {
				return nil, []error{err}
			}

			userByID := map[int64]*model.User{}
			for rows.Next() {
				var user model.User
				db.ScanRows(rows, &user)
				if user.ID == 0 {
					// no value returned
				} else {
					userByID[int64(user.ID)] = &user
				}
			}

			orderedUsers := make([]*model.User, len(ids))
			for i, id := range ids {
				orderedUsers[i] = userByID[id]
			}

			return orderedUsers, nil
		},
	})
}

func newUsersByCompanyIDLoader(db *gorm.DB) *generated.UserSliceLoader {
	return generated.NewUserSliceLoader(generated.UserSliceLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(companyIDs []int64) ([][]*model.User, []error) {
			rows, err := db.Model(&model.User{}).Where("company_id IN (?)", companyIDs).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			groupByCompanyID := make(map[int64][]*model.User, len(companyIDs))
			for rows.Next() {
				var user model.User
				db.ScanRows(rows, &user)
				if user.ID == 0 {
					// no value returned
				} else {
					groupByCompanyID[int64(user.CompanyID)] = append(groupByCompanyID[int64(user.CompanyID)], &user)
				}
			}

			orderedUsers := make([][]*model.User, len(companyIDs))
			for i, companyID := range companyIDs {
				orderedUsers[i] = groupByCompanyID[companyID]
			}

			return orderedUsers, nil
		},
	})
}
func newUserByEmailLoader(db *gorm.DB) *generated.UserStringLoader {
	return generated.NewUserStringLoader(generated.UserStringLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(emails []string) ([]*model.User, []error) {
			rows, err := db.Model(&model.User{}).Where("email IN (?)", emails).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			userByEmail := map[string]*model.User{}
			for rows.Next() {
				var user model.User
				db.ScanRows(rows, &user)
				if user.ID == 0 {
					// no value returned
				} else {
					userByEmail[user.Email] = &user
				}
			}

			orderedUsers := make([]*model.User, len(emails))
			for i, email := range emails {
				orderedUsers[i] = userByEmail[email]
			}

			return orderedUsers, nil
		},
	})
}
