package reset_password

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/server/entities/user"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/reset_password"
	"github.com/pkg/errors"
)

// ResetPassword has a reset password model object
type ResetPassword struct {
	ResetPasswordModel *reset_password.ResetPassword
	database           *sql.DB
}

// New returns a reset password entity
func New(id, userID int64, token string, expiresAt time.Time) *ResetPassword {
	return &ResetPassword{
		ResetPasswordModel: reset_password.New(id, userID, token, expiresAt),
		database:           db.SharedInstance().Connection,
	}
}

// GenerateResetPassword generate new token and return a new reset password entity
func GenerateResetPassword(userID int64, email string) (*ResetPassword, error) {
	// tokenを生成
	h := md5.New()
	_, err := io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return nil, errors.Wrap(err, "token generate error")
	}
	_, err = io.WriteString(h, email)
	if err != nil {
		return nil, errors.Wrap(err, "token generate error")
	}
	token := fmt.Sprintf("%x", h.Sum(nil))

	return New(0, userID, token, time.Now().AddDate(0, 0, 1)), nil
}

// FindAvailable find available reset password entity
func FindAvailable(id int64, token string) (*ResetPassword, error) {
	r, err := reset_password.FindAvailable(id, token)
	if err != nil {
		return nil, err
	}
	return &ResetPassword{
		ResetPasswordModel: r,
		database:           db.SharedInstance().Connection,
	}, nil
}

// Authenticate call authenticate in model
func Authenticate(id int64, token string) error {
	return reset_password.Authenticate(id, token)
}

// Save call save in model
func (r *ResetPassword) Save() error {
	return r.ResetPasswordModel.Save()
}

// User returns a owner user entity
func (r *ResetPassword) User() (*user.User, error) {
	var userID int64
	err := r.database.QueryRow("select user_id from reset_passwords where id = ?;", r.ResetPasswordModel.ID).Scan(&userID)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	u, err := user.Find(userID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ChangeUserPassword change password in owner user record
func (r *ResetPassword) ChangeUserPassword(password string) (*user.User, error) {
	u, err := r.User()
	if err != nil {
		return nil, err
	}

	hashPassword, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}

	tx, err := r.database.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction start error")
	}
	_, err = tx.Exec("update users set password = ? where id = ?;", hashPassword, r.ResetPasswordModel.ID)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql execute error")
	}

	_, err = tx.Exec("update reset_passwords set expires_at = now() where id = ?;", r.ResetPasswordModel.ID)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql execute error")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "transaction commit error")
	}
	u, err = r.User()
	if err != nil {
		return nil, err
	}
	return u, nil
}
