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


func TestRegistration(t *testing.T) {
	email := "registration@example.com"
	password := "hogehoge"

	reg := user.Registration(email, password)
	if reg != true {
		t.Fatalf("登録できない")
	}

	mydb := &db.Database{}
	var database db.DB = mydb
	sql := database.Init()
	rows, _ := sql.Query("select * from users where email = ?;", email)

	id, dbemail, dbpassword, created_at, updated_at := 0, "", "", "", ""
	for rows.Next() {
		_ = rows.Scan(&id, &dbemail, &dbpassword, &created_at, &updated_at)
	}
	if dbemail == "" {
		t.Error("ユーザが登録できていない")
	}

	reg = user.Registration(email, password)
	if reg != false {
		t.Error("ユーザが二重登録できている")
	}
}

func TestLogin(t *testing.T) {
	email := "login@example.com"
	password := "hogehoge"

	_ = user.Registration(email, password)

	current_user, err := user.Login(email, password)
	if err != nil {
		t.Error("ログイン時にエラー発生")
	}

	if current_user.Email != email {
		t.Error("ログインできない")
	}

	current_user, err = user.Login(email, "fugafuga")
	if current_user.Email == email {
		t.Error("パスワードが違うはずなのにログインできる")
	}

	current_user, err = user.Login("hogehoge@example.com", password)
	if current_user.Email == email {
		t.Error("メールアドレスが違うはずなのにログインできる")
	}

	current_user, err = user.Login("hogehoge@example.com", "fugafuga")
	if current_user.Email == email {
		t.Error("メールアドレスもパスワードも違うはずなのにログインできる")
	}
}
