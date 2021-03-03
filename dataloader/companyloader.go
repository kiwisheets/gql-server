package dataloader

import (
	"time"

	"github.com/kiwisheets/gql-server/dataloader/generated"
	"github.com/kiwisheets/gql-server/model"
	"gorm.io/gorm"
)

func newCompanyByIDLoader(db *gorm.DB) *generated.CompanyLoader {
	return generated.NewCompanyLoader(generated.CompanyLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []int64) ([]*model.Company, []error) {
			rows, err := db.Model(&model.Company{}).Where(ids).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			companyByID := map[int64]*model.Company{}
			for rows.Next() {
				var company model.Company
				db.ScanRows(rows, &company)
				if company.ID == 0 {
					// no value returned
				} else {
					companyByID[int64(company.ID)] = &company
				}
			}

			orderedCompanies := make([]*model.Company, len(ids))
			for i, id := range ids {
				orderedCompanies[i] = companyByID[id]
				i++
			}

			return orderedCompanies, nil
		},
	})
}

// TODO: Fix this, what is this
func newCompanyByUserIDLoader(db *gorm.DB) *generated.CompanyLoader {
	return generated.NewCompanyLoader(generated.CompanyLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(userIDs []int64) ([]*model.Company, []error) {
			// get users
			userRows, err := db.Model(&model.User{}).Select("id, company_id").Where(userIDs).Rows()

			if err != nil {
				if userRows == nil {
					return nil, []error{err}
				}
			}
			defer userRows.Close()

			// map user IDs to company IDs and company IDs to nil
			// map to nil so that we only get unique companies from DB

			companyIDbyUserID := map[int64]int64{}
			companyByCompanyID := map[int64]*model.Company{}

			for userRows.Next() {
				var user model.User
				db.ScanRows(userRows, &user)
				companyIDbyUserID[int64(user.ID)] = int64(user.CompanyID)
				companyByCompanyID[int64(user.CompanyID)] = nil
			}

			// convert map to slice of now unique company IDs

			companyIDs := []int64{}
			for id := range companyByCompanyID {
				companyIDs = append(companyIDs, id)
			}

			companyRows, err := db.Model(&model.Company{}).Where(companyIDs).Rows()

			if err != nil {
				if companyRows == nil {
					return nil, []error{err}
				}
			}
			defer companyRows.Close()

			for companyRows.Next() {
				var company model.Company
				db.ScanRows(companyRows, &company)
				if company.ID == 0 {
					// no value returned
				} else {
					companyByCompanyID[int64(company.ID)] = &company
				}
			}

			companyByUserID := map[int64]*model.Company{}
			for userID, companyID := range companyIDbyUserID {
				companyByUserID[userID] = companyByCompanyID[companyID]
			}

			orderedCompanies := make([]*model.Company, len(userIDs))
			for i, id := range userIDs {
				orderedCompanies[i] = companyByUserID[id]
			}

			return orderedCompanies, nil
		},
	})
}

func newCompanyByCodeLoader(db *gorm.DB) *generated.CompanyStringLoader {
	return generated.NewCompanyStringLoader(generated.CompanyStringLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(companyCodes []string) ([]*model.Company, []error) {
			rows, err := db.Model(&model.Company{}).Where("code IN (?)", companyCodes).Rows()

			if err != nil {
				if rows == nil {
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			companyByCode := map[string]*model.Company{}
			for rows.Next() {
				var company model.Company
				db.ScanRows(rows, &company)
				if company.ID == 0 {
					// no value returned
				} else {
					companyByCode[company.Code] = &company
				}
			}

			orderedCompanies := make([]*model.Company, len(companyCodes))
			for i, code := range companyCodes {
				orderedCompanies[i] = companyByCode[code]
			}

			return orderedCompanies, nil
		},
	})
}
