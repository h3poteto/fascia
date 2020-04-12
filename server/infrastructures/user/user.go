package user

import (
	"database/sql"

	"github.com/h3poteto/fascia/server/domains/user"
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
func (u *User) Find(id int64) (*user.User, error) {
	var uuid sql.NullInt64
	var email, password string
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := u.db.QueryRow("SELECT email, password, provider, oauth_token, user_name, uuid, avatar_url FROM users WHERE id = $1;", id).Scan(&email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "user repository")
	}
	return &user.User{
		ID:             id,
		Email:          email,
		HashedPassword: password,
		Provider:       provider,
		OauthToken:     oauthToken,
		UUID:           uuid,
		UserName:       userName,
		Avatar:         avatarURL,
	}, nil
}

// FindByEmail search a user according to email
func (u *User) FindByEmail(email string) (*user.User, error) {
	var id int64
	var password string
	var uuid sql.NullInt64
	var provider, oauthToken, userName, avatarURL sql.NullString
	err := u.db.QueryRow("SELECT id, email, password, provider, oauth_token, user_name, uuid, avatar_url FROM users WHERE email = $1;", email).Scan(&id, &email, &password, &provider, &oauthToken, &userName, &uuid, &avatarURL)
	if err != nil {
		return nil, errors.Wrap(err, "user repository")
	}
	return &user.User{
		ID:             id,
		Email:          email,
		HashedPassword: password,
		Provider:       provider,
		OauthToken:     oauthToken,
		UUID:           uuid,
		UserName:       userName,
		Avatar:         avatarURL,
	}, nil
}

// Save save user model in database
func (u *User) Create(email, hashedPassword string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) (int64, error) {
	var id int64
	err := u.db.QueryRow("INSERT INTO users (email, password, provider, oauth_token, uuid, user_name, avatar_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", email, hashedPassword, provider, oauthToken, uuid, userName, avatar).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "user repository")
	}
	return id, nil
}

// Update update user model in database
func (u *User) Update(id int64, email string, provider sql.NullString, oauthToken sql.NullString, uuid sql.NullInt64, userName sql.NullString, avatar sql.NullString) error {
	_, err := u.db.Exec("UPDATE users SET email = $1, provider = $2, oauth_token = $3, uuid = $4, user_name = $5, avatar_url = $6 WHERE id = $7;", email, provider, oauthToken, uuid, userName, avatar, id)
	if err != nil {
		return errors.Wrap(err, "user repository")
	}
	return nil
}

// UpdatePassword update only user password.
func (u *User) UpdatePassword(id int64, password string) error {
	_, err := u.db.Exec("UPDATE users SET password = $1 WHERE id = $2", password, id)
	if err != nil {
		return errors.Wrap(err, "user repository")
	}
	return nil
}
