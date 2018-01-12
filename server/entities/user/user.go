package user

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/infrastructures/user"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	infrastructure *user.User
}

// New returns a user entity
func New(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *User {
	infrastructure := user.New(id, email, provider, oauthToken, uuid, userName, avatar)
	u := &User{
		infrastructure: infrastructure,
	}
	u.reload()
	return u
}

func (u *User) reflect() {
	u.infrastructure.ID = u.ID
	u.infrastructure.Email = u.Email
	u.infrastructure.Provider = u.Provider
	u.infrastructure.OauthToken = u.OauthToken
	u.infrastructure.UUID = u.UUID
	u.infrastructure.UserName = u.UserName
	u.infrastructure.Avatar = u.Avatar
}

func (u *User) reload() error {
	if u.ID != 0 {
		latestUser, err := user.Find(u.ID)
		if err != nil {
			return err
		}
		u.infrastructure = latestUser
	}
	u.ID = u.infrastructure.ID
	u.Email = u.infrastructure.Email
	u.Provider = u.infrastructure.Provider
	u.OauthToken = u.infrastructure.OauthToken
	u.UUID = u.infrastructure.UUID
	u.UserName = u.infrastructure.UserName
	u.Avatar = u.infrastructure.Avatar
	return nil
}

// Registration create a new user record
func Registration(email, password, passwordConfirm string) (*User, error) {
	infrastructure, err := user.Registration(email, password, passwordConfirm)
	if err != nil {
		return nil, err
	}
	u := &User{
		infrastructure: infrastructure,
	}
	if err := u.reload(); err != nil {
		return nil, err
	}
	return u, nil
}

// Login authenticate email and password
func Login(userEmail string, userPassword string) (*User, error) {
	u, err := FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	bytePassword := []byte(userPassword)
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "password did not match")
	}
	return u, nil
}

// FindOrCreateFromGithub create or update user table base on github information
func FindOrCreateFromGithub(githubUser *github.User, token string, primaryEmail string) (*User, error) {
	u, err := FindByEmail(primaryEmail)
	if err != nil {
		// 見つからない場合は新規登録する
		if err := u.infrastructure.CreateGithubUser(token, githubUser, primaryEmail); err != nil {
			return nil, err
		}
	}

	// 登録されているOAuth情報が更新された場合には，合わせてレコードも更新しておく
	if !u.OauthToken.Valid || u.OauthToken.String != token {
		if err := u.infrastructure.UpdateGithubUserInfo(token, githubUser); err != nil {
			return nil, err
		}
	}

	return u, nil
}
