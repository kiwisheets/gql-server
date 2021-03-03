package model

type AuthData struct {
	User             *User   `json:"user"`
	Token            *string `json:"token"`
	TwoFactorEnabled bool    `json:"twoFactorEnabled"`
}
