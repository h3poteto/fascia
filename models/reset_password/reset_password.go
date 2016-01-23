package reset_password

import (
	"../../modules/logging"
	"../db"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

type ResetPassword interface {
}

type ResetPasswordStruct struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiresAt time.Time
	database  db.DB
}

func NewResetPassword(id int64, userId int64, token string, expiresAt time.Time) *ResetPasswordStruct {
	resetPassword := &ResetPasswordStruct{Id: id, UserId: userId, Token: token, ExpiresAt: expiresAt}
	resetPassword.Initialize()
	return resetPassword
}

func GenerateResetPassword(userId int64, email string) *ResetPasswordStruct {
	// tokenを生成
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, email)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return NewResetPassword(0, userId, token, time.Now())
}

func (u *ResetPasswordStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *ResetPasswordStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into reset_passwords (user_id, token, expires_at, created_at) values (?, ?, ?, now());", u.UserId, u.Token, u.ExpiresAt)
	if err != nil {
		logging.SharedInstance().MethodInfo("ResetPassword", "Save").Errorf("reset_password save error: %v", err)
		return false
	}
	u.Id, _ = result.LastInsertId()
	return true
}
