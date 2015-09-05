package user

import (
	"os"
	"testing"
	"../models/user"
	"../models/db"
)

func TestMain(m *testing.M) {
	// TODO: ここあとで共通化したい
	testdb := os.Getenv("DB_TEST_NAME")
	currentdb := os.Getenv("DB_NAME")
	os.Setenv("DB_NAME", testdb)

	code := m.Run()
	mydb := &db.Database{}
	var database db.DB = mydb
	sql := database.Init()
	sql.Exec("truncate table users;")
	sql.Close()
	os.Setenv("DB_NAME", currentdb)
	os.Exit(code)
}

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
