package reset_password

import (
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/user"

	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type ResetPassword interface {
}

type ResetPasswordStruct struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	database  *sql.DB
}

func NewResetPassword(id int64, userID int64, token string, expiresAt time.Time) *ResetPasswordStruct {
	resetPassword := &ResetPasswordStruct{ID: id, UserID: userID, Token: token, ExpiresAt: expiresAt}
	resetPassword.Initialize()
	return resetPassword
}

func GenerateResetPassword(userID int64, email string) *ResetPasswordStruct {
	// tokenを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, email)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return NewResetPassword(0, userID, token, time.Now().AddDate(0, 0, 1))
}

func Authenticate(id int64, token string) error {
	database := db.SharedInstance().Connection

	var targetID int64
	err := database.QueryRow("select id from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&targetID)
	if err != nil {
		return errors.Wrap(err, "sql select error")
	}

	return nil
}

func ChangeUserPassword(id int64, token string, password string) (u *user.UserStruct, e error) {
	database := db.SharedInstance().Connection
	tx, _ := database.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			u = nil
			switch ty := err.(type) {
			case runtime.Error:
				e = errors.Wrap(ty, "runtime error")
			case string:
				e = errors.New(err.(string))
			default:
				e = errors.New("unexpected error")
			}
		}
	}()

	var userID int64
	err := tx.QueryRow("select user_id from reset_passwords where id = ? and token = ? and expires_at > now();", id, token).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql select error")
	}

	hashPassword, err := user.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("update users set password = ? where id = ?;", hashPassword, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql execute error")
	}

	_, err = tx.Exec("update reset_passwords set expires_at = now() where id = ?;", id)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql execute error")
	}

	u, err = user.FindUser(userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return u, nil
}

func (u *ResetPasswordStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *ResetPasswordStruct) Save() error {
	result, err := u.database.Exec("insert into reset_passwords (user_id, token, expires_at, created_at) values (?, ?, ?, now());", u.UserID, u.Token, u.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}
