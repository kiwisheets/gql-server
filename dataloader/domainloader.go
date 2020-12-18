package dataloader

import (
	"time"

	"github.com/kiwisheets/gql-server/dataloader/generated"
	"github.com/kiwisheets/gql-server/orm/model"
	"gorm.io/gorm"
)

func newDomainsByCompanyIDLoader(db *gorm.DB) *generated.DomainSliceLoader {
	return generated.NewDomainSliceLoader(generated.DomainSliceLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(companyIDs []int64) ([][]*model.Domain, []error) {
			rows, err := db.Model(&model.Domain{}).Where("company_id IN (?)", companyIDs).Rows()

			if err != nil {
				if rows != nil {
					rows.Close()
					return nil, []error{err}
				}
				// log error
			}
			defer rows.Close()

			groupByCompanyID := make(map[int64][]*model.Domain, len(companyIDs))
			for rows.Next() {
				var domain model.Domain
				db.ScanRows(rows, &domain)
				if domain.ID == 0 {
					// no value returned
				} else {
					groupByCompanyID[int64(domain.CompanyID)] = append(groupByCompanyID[int64(domain.CompanyID)], &domain)
				}
			}

			companyDomains := make([][]*model.Domain, len(companyIDs))
			var errs []error
			for i, companyID := range companyIDs {
				companyDomains[i] = groupByCompanyID[companyID]
			}

			return companyDomains, errs
		},
	})
}
