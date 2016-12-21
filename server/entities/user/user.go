package user

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/user"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User has a user model object
type User struct {
	UserModel *user.User
	database  *sql.DB
}

// New returns a user entity
func New(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *User {
	return &User{
		UserModel: user.New(id, email, provider, oauthToken, uuid, userName, avatar),
		database:  db.SharedInstance().Connection,
	}
}

// Registration create a new user record
func Registration(email, password, passwordConfirm string) (*User, error) {
	u, err := user.Registration(email, password, passwordConfirm)
	if err != nil {
		return nil, err
	}
	return &User{
		UserModel: u,
		database:  db.SharedInstance().Connection,
	}, nil
}

// Find returns a user entity
func Find(id int64) (*User, error) {
	u, err := user.Find(id)
	if err != nil {
		return nil, err
	}
	return &User{
		UserModel: u,
		database:  db.SharedInstance().Connection,
	}, nil
}

// FindByEmail returns a user entity
func FindByEmail(email string) (*User, error) {
	u, err := user.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return &User{
		UserModel: u,
		database:  db.SharedInstance().Connection,
	}, nil
}

// Login authenticate email and password
func Login(userEmail string, userPassword string) (*User, error) {
	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", userEmail).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}

	u := New(id, email, provider, oauthToken, uuid, userName, avatarURL)
	bytePassword := []byte(userPassword)
	err = bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "password did not match")
	}
	return u, nil
}

// FindOrCreateFromGithub create or update user table base on github information
func FindOrCreateFromGithub(githubUser *github.User, token string, primaryEmail string) (*User, error) {
	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	rows, err := database.Query("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where uuid = ? or email = ?;", *githubUser.ID, primaryEmail)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	// 新規登録の場合には見つからないのでエラーにならないようにscanする
	for rows.Next() {
		err := rows.Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
		if err != nil {
			return nil, errors.Wrap(err, "sql scan error")
		}
	}
	u := New(id, email, provider, oauthToken, uuid, userName, avatarURL)

	// id==0, 即ち初期値の場合には新規登録する
	// TODO: できればrowsの長さだけで判定したい
	if id == 0 {
		if err := u.UserModel.CreateGithubUser(token, githubUser, primaryEmail); err != nil {
			return nil, err
		}
	}

	// 登録されているOAuth情報が更新された場合には，合わせてレコードも更新しておく
	if !u.UserModel.OauthToken.Valid || u.UserModel.OauthToken.String != token {
		if err := u.UserModel.UpdateGithubUserInfo(token, githubUser); err != nil {
			return nil, err
		}
	}

	return u, nil
}

// HashPassword returns a password
func HashPassword(password string) ([]byte, error) {
	return user.HashPassword(password)
}

// Projects list up projects related a user
func (u *User) Projects() ([]*project.Project, error) {
	var slice []*project.Project
	rows, err := u.database.Query("select id, user_id, repository_id, title, description, show_issues, show_pull_requests from projects where user_id = ?;", u.UserModel.ID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	for rows.Next() {
		var id, userID int64
		var repositoryID sql.NullInt64
		var title string
		var description string
		var showIssues, showPullRequests bool
		err := rows.Scan(&id, &userID, &repositoryID, &title, &description, &showIssues, &showPullRequests)
		if err != nil {
			return nil, errors.Wrap(err, "sql select error")
		}
		if id != 0 {
			p := project.New(id, userID, title, description, repositoryID, showIssues, showPullRequests)
			slice = append(slice, p)
		}
	}
	return slice, nil
}
