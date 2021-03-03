// Package generated contains generated dataloader configurations
package generated

//go:generate go run github.com/vektah/dataloaden UserLoader int64 *github.com/kiwisheets/gql-server/model.User
//go:generate go run github.com/vektah/dataloaden UserSliceLoader int64 []*github.com/kiwisheets/gql-server/model.User
//go:generate go run github.com/vektah/dataloaden UserStringLoader string *github.com/kiwisheets/gql-server/model.User
//go:generate go run github.com/vektah/dataloaden UserStringSliceLoader string []*github.com/kiwisheets/gql-server/model.User

//go:generate go run github.com/vektah/dataloaden CompanyLoader int64 *github.com/kiwisheets/gql-server/model.Company
//go:generate go run github.com/vektah/dataloaden CompanyStringLoader string *github.com/kiwisheets/gql-server/model.Company

//go:generate go run github.com/vektah/dataloaden DomainSliceLoader int64 []*github.com/kiwisheets/gql-server/model.Domain

//ignore me go:generate go run github.com/vektah/dataloaden PermissionsLoader int64 []*github.com/kiwisheets/auth/permission.Permission
//go:generate go run github.com/vektah/dataloaden PermissionsLoader int64 []*github.com/kiwisheets/auth/permission.Permission

// Checks if user has permission by ID
//go:generate go run github.com/vektah/dataloaden RoleLoader int64 []github.com/kiwisheets/auth/permission.Role

//go:generate go run github.com/vektah/dataloaden AddressLoader int64 *github.com/kiwisheets/gql-server/model.Address
