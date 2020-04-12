package views

import (
	"github.com/h3poteto/fascia/server/domains/user"
)

// User is a type of user entity.
type User struct {
	ID       int64  `json:"ID"`
	Email    string `json:"Email"`
	UserName string `json:"UserName"`
	Avatar   string `json:"Avatar"`
}

// ParseUserJSON returns json formatted user entity.
func ParseUserJSON(user *user.User) (*User, error) {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		UserName: user.UserName.String,
		Avatar:   user.Avatar.String,
	}, nil
}
