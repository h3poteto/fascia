package user

import (
	"database/sql"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/server/domains/project"
)

// User has a user model object
type User struct {
	ID             int64
	Email          string
	HashedPassword string
	Provider       sql.NullString
	OauthToken     sql.NullString
	UUID           sql.NullInt64
	UserName       sql.NullString
	Avatar         sql.NullString
}

// New returns a User struct.
func New(id int64, email, hashedPassword string, provider, oauthToken sql.NullString, uuid sql.NullInt64, userName, avatar sql.NullString) *User {
	return &User{
		id,
		email,
		hashedPassword,
		provider,
		oauthToken,
		uuid,
		userName,
		avatar,
	}
}

// Update updates a user entity.
func (u *User) Update(email, hashedPassword string, provider, oauthToken sql.NullString, uuid sql.NullInt64, userName, avatar sql.NullString) {
	u.Email = email
	u.HashedPassword = hashedPassword
	u.Provider = provider
	u.OauthToken = oauthToken
	u.UUID = uuid
	u.UserName = userName
	u.Avatar = avatar
}

// UpdateGithubUser updates a user entity with github user.
func (u *User) UpdateGithubUser(githubUser *github.User, id int64, email, token string) {
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}
	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.UUID = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
}

// Projects list up projects related a user
func (u *User) Projects(infrastructure project.Repository) ([]*project.Project, error) {
	return project.Projects(u.ID, infrastructure)
}
