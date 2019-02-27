package reset_password

import (
	"github.com/h3poteto/fascia/lib/modules/database"

	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// ResetPassword has reset password record
type ResetPassword struct {
	db *sql.DB
}

// New returns a reset password
func New(db *sql.DB) *ResetPassword {
	resetPassword := &ResetPassword{
		db,
	}
	return resetPassword
}

// Authenticate check token with record
func (r *ResetPassword) Authenticate(id int64, token string) error {
	var targetID int64
	err := r.db.QueryRow("select id from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&targetID)
	if err != nil {
		return errors.Wrap(err, "authenticate error")
	}

	return nil
}

// FindAvailable search a reset password which can authorize
func (r *ResetPassword) FindAvailable(id int64, token string) (int64, int64, string, time.Time, error) {
	var userID int64
	var expiresAt time.Time
	err := r.db.QueryRow("select user_id, expires_at from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&userID, &expiresAt)
	return id, userID, token, expiresAt, err
}

// Find find a reset password by id.
func (r *ResetPassword) Find(id int64) (int64, int64, string, time.Time, error) {
	var userID int64
	var token string
	var expiresAt time.Time
	db := database.SharedInstance().Connection
	err := db.QueryRow("select user_id, token, expires_at from reset_passwords where id = ?;", id).Scan(&userID, &token, &expiresAt)
	return id, userID, token, expiresAt, err
}

// Create save object to record
func (r *ResetPassword) Create(userID int64, token string, expiresAt time.Time) (int64, error) {
	result, err := r.db.Exec("insert into reset_passwords (user_id, token, expires_at, created_at) values (?, ?, ?, now());", userID, token, expiresAt)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// UpdateExpire update expires to now.
func (r *ResetPassword) UpdateExpire(id int64) error {
	_, err := r.db.Exec("update reset_passwords set expires_at = now() where id = ?;", id)
	if err != nil {
		return err
	}
	return nil
}
