package services

import (
	"github.com/h3poteto/fascia/server/aggregations/user"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type User struct {
	UserAggregation *user.User
}

func NewUserService(userAg *user.User) *User {
	return &User{
		UserAggregation: userAg,
	}
}

func CurrentUser(userID int64) (*User, error) {
	userAggregation, err := user.CurrentUser(userID)
	if err != nil {
		return nil, err
	}
	return &User{
		UserAggregation: userAggregation,
	}, nil
}

func LoginUser(email, password string) (*User, error) {
	userAggregation, err := user.Login(email, password)
	if err != nil {
		return nil, err
	}
	return &User{
		UserAggregation: userAggregation,
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

	userAggregation, err := user.FindOrCreateFromGithub(githubUser, token, primaryEmail)
	if err != nil {
		return nil, err
	}
	return &User{
		UserAggregation: userAggregation,
	}, nil
}

func (u *User) Projects() ([]*Project, error) {
	projectAggregations, err := u.UserAggregation.Projects()
	if err != nil {
		return nil, err
	}
	var slice []*Project
	for _, p := range projectAggregations {
		slice = append(slice, NewProjectService(p))
	}
	return slice, nil
}
