package reset_password

import (
	"github.com/h3poteto/fascia/lib/modules/database"

	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// ResetPassword has reset password record
type ResetPassword struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	db        *sql.DB
}

// New returns a reset password
func New(id int64, userID int64, token string, expiresAt time.Time) *ResetPassword {
	resetPassword := &ResetPassword{ID: id, UserID: userID, Token: token, ExpiresAt: expiresAt}
	resetPassword.initialize()
	return resetPassword
}

// Authenticate check token with record
func Authenticate(id int64, token string) error {
	db := database.SharedInstance().Connection

	var targetID int64
	err := db.QueryRow("select id from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&targetID)
	if err != nil {
		return errors.Wrap(err, "authenticate error")
	}

	return nil
}

// FindAvailable search a reset password which can authorize
func FindAvailable(id int64, token string) (*ResetPassword, error) {
	var userID int64
	var expiresAt time.Time
	db := database.SharedInstance().Connection
	err := db.QueryRow("select user_id, expires_at from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&userID, &expiresAt)
	if err != nil {
		return nil, errors.Wrap(err, "find available error")
	}
	return New(id, userID, token, expiresAt), nil
}

// Find find a reset password by id.
func Find(id int64) (*ResetPassword, error) {
	var userID int64
	var token string
	var expiresAt time.Time
	db := database.SharedInstance().Connection
	err := db.QueryRow("select user_id, token, expires_at from reset_passwords where id = ?;", id).Scan(&userID, &token, &expiresAt)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, userID, token, expiresAt), nil
}

func (r *ResetPassword) initialize() {
	r.db = database.SharedInstance().Connection
}

// Save save object to record
func (r *ResetPassword) Save() error {
	result, err := r.db.Exec("insert into reset_passwords (user_id, token, expires_at, created_at) values (?, ?, ?, now());", r.UserID, r.Token, r.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "save error")
	}
	r.ID, _ = result.LastInsertId()
	return nil
}

// UpdateExpire update expires to now.
func (r *ResetPassword) UpdateExpire(tx *sql.Tx) error {
	if tx != nil {
		_, err := tx.Exec("update reset_passwords set expires_at = now() where id = ?;", r.ID)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "sql execute error")
		}
	} else {
		_, err := r.db.Exec("update reset_passwords set expires_at = now() where id = ?;", r.ID)
		if err != nil {
			return errors.Wrap(err, "sql execute error")
		}
	}
	return nil
}
