package reset_password

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/entities/user"
	"github.com/h3poteto/fascia/server/infrastructures/reset_password"
	"github.com/pkg/errors"
)

// ResetPassword has a reset password model object
type ResetPassword struct {
	ID             int64
	UserID         int64
	Token          string
	ExpiresAt      time.Time
	infrastructure *reset_password.ResetPassword
}

// New returns a reset password entity
func New(id, userID int64, token string, expiresAt time.Time) *ResetPassword {
	infrastructure := reset_password.New(id, userID, token, expiresAt)
	r := &ResetPassword{
		infrastructure: infrastructure,
	}
	r.reload()
	return r
}

func (r *ResetPassword) reflect() {
	r.infrastructure.ID = r.ID
	r.infrastructure.UserID = r.UserID
	r.infrastructure.Token = r.Token
	r.infrastructure.ExpiresAt = r.ExpiresAt
}

func (r *ResetPassword) reload() error {
	if r.ID != 0 {
		latestReset, err := reset_password.Find(r.ID)
		if err != nil {
			return err
		}
		r.infrastructure = latestReset
	}
	r.ID = r.infrastructure.ID
	r.UserID = r.infrastructure.UserID
	r.Token = r.infrastructure.Token
	r.ExpiresAt = r.infrastructure.ExpiresAt
	return nil
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

// Authenticate call authenticate in model
func Authenticate(id int64, token string) error {
	return reset_password.Authenticate(id, token)
}

// Save call save in model
func (r *ResetPassword) Save() error {
	r.reflect()
	return r.infrastructure.Save()
}

// ChangeUserPassword change password in owner user record
func (r *ResetPassword) ChangeUserPassword(password string) (*user.User, error) {
	u, err := r.User()
	if err != nil {
		return nil, err
	}

	// Control transaction in here.
	db := database.SharedInstance().Connection
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "transaction start error")
	}

	if err := u.UpdatePassword(password, tx); err != nil {
		return nil, err
	}

	if err := r.infrastructure.UpdateExpire(tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "transaction commit error")
	}
	return r.User()
}
