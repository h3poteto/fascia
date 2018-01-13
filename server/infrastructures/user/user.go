package user

import (
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/lib/modules/logging"

	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User has a user record
type User struct {
	ID         int64
	Email      string
	Password   string
	Provider   sql.NullString
	OauthToken sql.NullString
	UUID       sql.NullInt64
	UserName   sql.NullString
	Avatar     sql.NullString
	db         *sql.DB
}

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

// hashPassword generate hash password
func hashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashed, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate password")
	}
	err = bcrypt.CompareHashAndPassword(hashed, bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "did not match password")
	}
	return hashed, nil
}

// New returns a user object
func New(id int64, email string, password string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) *User {
	user := &User{ID: id, Email: email, Password: password, Provider: provider, OauthToken: oauthToken, UUID: uuid, UserName: userName, Avatar: avatar}
	user.initialize()
	return user
}

func (u *User) initialize() {
	u.db = database.SharedInstance().Connection
}

// Registration is create new user through validation
func Registration(email string, password string, passwordConfirm string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := New(0, email, string(hashed), sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	err = user.Save()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Find search a user according to id
func Find(id int64) (*User, error) {
	db := database.SharedInstance().Connection

	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := db.QueryRow("select email, password, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", id).Scan(&email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatarURL), nil
}

// FindByEmail search a user according to email
func FindByEmail(email string) (*User, error) {
	db := database.SharedInstance().Connection
	var id int64
	var password string
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := db.QueryRow("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", email).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatarURL), nil
}

// Save save user model in database
func (u *User) Save() error {
	// TODO: この前にvalidationを入れたい
	result, err := u.db.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", u.Email, u.Password, u.Provider, u.OauthToken, u.UUID, u.UserName, u.Avatar)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	logging.SharedInstance().MethodInfo("user", "Save").Infof("user saved: %v", u.ID)
	return nil
}

// Update update user model in database
func (u *User) Update() error {
	_, err := u.db.Exec("update users set provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where email = ?;", u.Provider, u.OauthToken, u.UUID, u.UserName, u.Avatar, u.Email)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	return nil
}

// CreateGithubUser create a user from github authentication
func CreateGithubUser(token string, githubUser *github.User, primaryEmail string) (*User, error) {
	bytePassword, err := hashPassword(randomString())
	if err != nil {
		return nil, err
	}
	provider := sql.NullString{String: "github", Valid: true}
	oauthToken := sql.NullString{String: token, Valid: true}
	userName := sql.NullString{String: *githubUser.Login, Valid: true}
	uuid := sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	avatar := sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	user := New(0, primaryEmail, string(bytePassword), provider, oauthToken, uuid, userName, avatar)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateGithubUserInfo update a user from github authentication
func (u *User) UpdateGithubUserInfo(token string, githubUser *github.User) error {
	u.Provider = sql.NullString{String: "github", Valid: true}
	u.OauthToken = sql.NullString{String: token, Valid: true}
	u.UserName = sql.NullString{String: *githubUser.Login, Valid: true}
	u.UUID = sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	u.Avatar = sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	if err := u.Update(); err != nil {
		return err
	}
	return nil
}

// UpdatePassword update password in user.
func (u *User) UpdatePassword(tx *sql.Tx) error {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	if tx != nil {
		_, err := tx.Exec("update users set password = ? where id = ?;", hashed, u.ID)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "sql execute error")
		}
	} else {
		_, err := u.db.Exec("update users set password = ? where id = ?;", hashed, u.ID)
		if err != nil {
			return errors.Wrap(err, "sql execute error")
		}
	}
	return nil
}
