package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/auth"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/dataloader"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/generated"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/graphql/modelgen"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/emvi/hide"
)

func (r *companyResolver) Users(ctx context.Context, obj *model.Company) ([]*model.User, error) {
	return dataloader.For(ctx).UsersByCompanyID.Load(obj.IDint())
}

func (r *companyResolver) Domains(ctx context.Context, obj *model.Company) ([]string, error) {
	domains, errs := dataloader.For(ctx).DomainsByCompanyID.Load(obj.IDint())

	domainStrings := make([]string, len(domains))
	for i, domain := range domains {
		domainStrings[i] = domain.Domain
	}

	return domainStrings, errs
}

func (r *mutationResolver) CreateCompany(ctx context.Context, company modelgen.CreateCompanyInput) (*model.Company, error) {
	companyObject := model.Company{
		Code: company.Code,
		Name: company.Name,
	}

	// domain strings to domain models
	for _, d := range company.Domains {
		companyObject.Domains = append(companyObject.Domains, model.Domain{
			Domain: d,
		})
	}

	if err := r.DB.Create(&companyObject).Error; err != nil {
		return nil, fmt.Errorf("Unable to create Company. Already exists")
	}

	return &companyObject, nil
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id hide.ID) (*bool, error) {
	err := r.DB.Delete(&model.Company{
		SoftDelete: model.SoftDelete{
			ID: id,
		},
	}).Error
	if err == nil {
		return util.Bool(true), nil
	}
	return util.Bool(false), err
}

func (r *queryResolver) CompanyName(ctx context.Context, code string) (*string, error) {
	company, err := dataloader.For(ctx).CompanyByCode.Load(code)

	if company == nil {
		return util.String(""), fmt.Errorf("No company exists")
	}

	return &company.Name, err
}

func (r *queryResolver) Company(ctx context.Context) (*model.Company, error) {
	return dataloader.For(ctx).CompanyByID.Load(auth.For(ctx).User.IDint())
}

func (r *queryResolver) OtherCompany(ctx context.Context, id hide.ID) (*model.Company, error) {
	return dataloader.For(ctx).CompanyByID.Load(int64(id))
}

func (r *queryResolver) Companies(ctx context.Context, page *int) ([]*model.Company, error) {
	limit := 20
	companies := make([]*model.Company, limit)
	if page == nil {
		page = util.Int(0)
	}
	r.DB.Order("name").Limit(limit).Offset(limit * *page).Find(&companies)

	return companies, nil
}

// Company returns generated.CompanyResolver implementation.
func (r *Resolver) Company() generated.CompanyResolver { return &companyResolver{r} }

type companyResolver struct{ *Resolver }
