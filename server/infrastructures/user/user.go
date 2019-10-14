package user

import (
	"database/sql"

	"github.com/pkg/errors"
)

// User has a user record
type User struct {
	db *sql.DB
}

// New returns a user object
func New(db *sql.DB) *User {
	return &User{
		db,
	}
}

// Find finds a user which has specified id.
func (u *User) Find(id int64) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error) {
	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := u.db.QueryRow("select email, password, provider, oauth_token, user_name, uuid, avatar_url from users where id = ?;", id).Scan(&email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return 0, "", "", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{}, errors.Wrap(err, "user repository")
	}
	return id, email, password, provider, oauthToken, uuid, userName, avatarURL, nil
}

// FindByEmail search a user according to email
func (u *User) FindByEmail(email string) (int64, string, string, sql.NullString, sql.NullString, sql.NullInt64, sql.NullString, sql.NullString, error) {
	var id int64
	var password string
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := u.db.QueryRow("select id, email, password, provider, oauth_token, user_name, uuid, avatar_url from users where email = ?;", email).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return 0, "", "", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{}, errors.Wrap(err, "user repository")
	}
	return id, email, password, provider, oauthToken, uuid, userName, avatarURL, nil
}

// Save save user model in database
func (u *User) Create(email, password string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) (int64, error) {
	result, err := u.db.Exec("insert into users (email, password, provider, oauth_token, uuid, user_name, avatar_url, created_at) values (?, ?, ?, ?, ?, ?, ?, now());", email, password, provider, oauthToken, uuid, userName, avatar)
	if err != nil {
		return 0, errors.Wrap(err, "user repository")
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// Update update user model in database
func (u *User) Update(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) error {
	_, err := u.db.Exec("update users set email = ?, provider = ?, oauth_token = ?, uuid = ?, user_name = ?, avatar_url = ? where id = ?;", email, provider, oauthToken, uuid, userName, avatar, id)
	return errors.Wrap(err, "user repository")
}
