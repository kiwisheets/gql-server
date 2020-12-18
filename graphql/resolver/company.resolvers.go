package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/emvi/hide"
	"github.com/kiwisheets/auth"
	"github.com/kiwisheets/gql-server/dataloader"
	"github.com/kiwisheets/gql-server/graphql/generated"
	"github.com/kiwisheets/gql-server/graphql/modelgen"
	"github.com/kiwisheets/gql-server/orm/model"
	"github.com/kiwisheets/gql-server/util"
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
	var company model.Company
	err := r.DB.Where("code = ?", code).First(&company).Error
	if err != nil {
		return nil, fmt.Errorf("No company exists")
	}

	return &company.Name, err
}

func (r *queryResolver) Company(ctx context.Context) (*model.Company, error) {
	var company model.Company
	err := r.DB.Where(auth.For(ctx).CompanyID).First(&company).Error
	return &company, err
}

func (r *queryResolver) OtherCompany(ctx context.Context, id hide.ID) (*model.Company, error) {
	var company model.Company
	err := r.DB.Where(id).First(&company).Error
	return &company, err
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
