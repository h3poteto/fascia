package reset_password

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/server/aggregations/user"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/models/reset_password"
	"github.com/pkg/errors"
)

type ResetPassword struct {
	ResetPasswordModel *reset_password.ResetPassword
	database           *sql.DB
}

func New(id, userID int64, token string, expiresAt time.Time) *ResetPassword {
	return &ResetPassword{
		ResetPasswordModel: reset_password.New(id, userID, token, expiresAt),
		database:           db.SharedInstance().Connection,
	}
}

func GenerateResetPassword(userID int64, email string) *ResetPassword {
	// tokenを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, email)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return New(0, userID, token, time.Now().AddDate(0, 0, 1))
}

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
