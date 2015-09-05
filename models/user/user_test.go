package user

import (
	"testing"
	"."
)

func TestRegistrationAndLogin(t *testing.T) {
	aUser := user.NewUser()
	var u user.User = aUser

	email := "sample@example.com"
	password := "hogehoge"

	reg := u.Registration(email, password)
	if reg != true {
		t.Fatalf("登録できない")
	}

	current_user, err := u.Login(email, password)
	if err != nil {
		t.Error("登録後ログインできない")
	}

	if current_user.Email != email {
		t.Error("登録ユーザが見つからない")
	}
}
