package services

import (
	"github.com/h3poteto/fascia/server/entities/user"

	"context"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"time"
)

// User has a user entity
type User struct {
	UserEntity *user.User
}

// NewUser returns a user service
func NewUser(entity *user.User) *User {
	return &User{
		UserEntity: entity,
	}
}

// RegistrationUser create a user with email, and password
func RegistrationUser(email, password, passwordConfirm string) (*User, error) {
	u, err := user.Registration(email, password, passwordConfirm)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: u,
	}, nil
}

// FindUser search a user
func FindUser(id int64) (*User, error) {
	entity, err := user.Find(id)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

// FindUserByEmail search a user according to email
func FindUserByEmail(email string) (*User, error) {
	entity, err := user.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

// LoginUser authenticate email and password
func LoginUser(email, password string) (*User, error) {
	entity, err := user.Login(email, password)
	if err != nil {
		return nil, err
	}
	return &User{
		UserEntity: entity,
	}, nil
}

// FindOrCreateUserFromGithub authenticate with github, and if user already exist returns a user entity
// If user does not exist, create a user and return a new user service
func FindOrCreateUserFromGithub(token string) (*User, error) {
	// github認証
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.Wrap(err, "github api error")
	}

	// TODO: primaryじゃないEmailも保存しておいてログインブロックに使いたい
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	emails, _, _ := client.Users.ListEmails(ctx, nil)
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

// Projects returns a related projects
func (u *User) Projects() ([]*Project, error) {
	projectEntities, err := u.UserEntity.Projects()
	if err != nil {
		return nil, err
	}
	var slice []*Project
	for _, p := range projectEntities {
		slice = append(slice, NewProject(p))
	}
	return slice, nil
}
