package services

import (
	"github.com/h3poteto/fascia/server/entities/user"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type User struct {
	UserEntity *user.User
}

func NewUserService(entity *user.User) *User {
	return &User{
		UserEntity: entity,
	}
}

func RegistrationUser(email, password, passwordConfirm string) (*User, error) {
	u, err := user.Registration(email, password, passwordConfirm)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: u,
	}, nil
}

func FindUser(id int64) (*User, error) {
	entity, err := user.Find(id)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

func FindUserByEmail(email string) (*User, error) {
	entity, err := user.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

func LoginUser(email, password string) (*User, error) {
	entity, err := user.Login(email, password)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

func FindOrCreateUserFromGithub(token string) (*User, error) {
	// github認証
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	githubUser, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.Wrap(err, "github api error")
	}

	// TODO: primaryじゃないEmailも保存しておいてログインブロックに使いたい
	emails, _, _ := client.Users.ListEmails(nil)
	var primaryEmail string
	for _, email := range emails {
		if *email.Primary {
			primaryEmail = *email.Email
		}
	}

	entity, err := user.FindOrCreateFromGithub(githubUser, token, primaryEmail)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

func (u *User) Projects() ([]*Project, error) {
	projectEntities, err := u.UserEntity.Projects()
	if err != nil {
		return nil, err
	}
	var slice []*Project
	for _, p := range projectEntities {
		slice = append(slice, NewProjectService(p))
	}
	return slice, nil
}
