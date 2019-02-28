package user

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/entities/project"
)

// User has a user model object
type User struct {
	ID             int64
	Email          string
	Password       string
	Provider       sql.NullString
	OauthToken     sql.NullString
	UUID           sql.NullInt64
	UserName       sql.NullString
	Avatar         sql.NullString
	infrastructure Repository
}

type Repository interface {
	Find(id int64) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error)
	FindByEmail(string) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error)
	Create(string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString) (int64, error)
	Update(int64, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString) error
	UpdatePassword(int64, string) error
}

func New(id int64, email, password string, provider, oauthToken sql.NullString, uuid sql.NullInt64, userName, avatar sql.NullString, infrastructure Repository) *User {
	return &User{
		id,
		email,
		password,
		provider,
		oauthToken,
		uuid,
		userName,
		avatar,
		infrastructure,
	}
}

// create crates a new record.
// It contains password and does not transform the password in this method.
func (u *User) create() error {
	id, err := u.infrastructure.Create(u.Email, u.Password, u.Provider, u.OauthToken, u.UUID, u.UserName, u.Avatar)
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

// Update updates a record.
// It does not contain password field. Please use UpdatePassword when you want to update password.
func (u *User) Update() error {
	return u.infrastructure.Update(u.ID, u.Email, u.Provider, u.OauthToken, u.UUID, u.UserName, u.Avatar)
}

// UpdatePassword updates user password.
func (u *User) UpdatePassword(password string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}
	if err := u.infrastructure.UpdatePassword(u.ID, string(hashed)); err != nil {
		return err
	}
	u.Password = password
	return nil
}

// Projects list up projects related a user
func (u *User) Projects() ([]*project.Project, error) {
	return project.Projects(u.ID)
}
