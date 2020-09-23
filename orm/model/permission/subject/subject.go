//go:generate go-enum -f=$GOFILE --marshal

package subject

// Subject is an enumeration of permission subject values
/*
ENUM(
None
Any

Me

User
Users
UserContact
OtherUser
OtherUsers

Company
OtherCompany
OtherCompanies

Client
Clients
ClientContact

Contact
Contacts
)
*/
type Subject int64
