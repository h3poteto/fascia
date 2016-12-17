package user

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/models/db"

	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int64
	Email      string
	Password   string
	Provider   sql.NullString
	OauthToken sql.NullString
	Uuid       sql.NullInt64
	UserName   sql.NullString
	Avatar     sql.NullString
	database   *sql.DB
}

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func HashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate password")
	}
	err = bcrypt.CompareHashAndPassword(hashPassword, bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "did not match password")
	}
	return hashPassword, nil
}

func New(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *User {
	user := &User{ID: id, Email: email, Provider: provider, OauthToken: oauthToken, Uuid: uuid, UserName: userName, Avatar: avatar}
	user.initialize()
	return user
}

func (u *User) initialize() {
	u.database = db.SharedInstance().Connection
}

// Registration is create new user through validation
func Registration(email string, password string, passwordConfirm string) (int64, error) {
	hashPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := New(0, email, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	user.Password = string(hashPassword)
	err = user.Save()

	return user.ID, err
}

func Find(id int64) (*User, error) {
	database := db.SharedInstance().Connection

	var uuid sql.NullInt64
	var email string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select email, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", id).Scan(&email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

func FindByEmail(email string) (*User, error) {
	database := db.SharedInstance().Connection
	var id int64
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := database.QueryRow("select id, email, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", email).Scan(&id, &email, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, email, provider, oauthToken, uuid, userName, avatarURL), nil
}

func (u *User) Save() error {
	// TODO: この前にvalidationを入れたい
	result, err := u.database.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save").Infof("user saved: %v", u.ID)
	return nil
}

func (u *User) Update() error {
	_, err := u.database.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.Uuid, u.UserName, u.Avatar, u.Email)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	return nil
}

func (u *User) CreateGithubUser(token string, githubUser *github.User, primaryEmail string) error {
	u.Email = primaryEmail
	bytePassword, err := HashPassword(randomString())
	if err != nil {
		return err
	}
	u.Password = string(bytePassword)
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}

	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	if err := u.Save(); err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateGithubUserInfo(token string, githubUser *github.User) error {
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}
	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.Uuid = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	if err := u.Update(); err != nil {
		return err
	}
	return nil
}
