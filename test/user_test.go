package user

import (
	"os"
	"testing"
	"database/sql"
	"../models/user"
	"../models/project"
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
	table := database.Init()
	rows, _ := table.Query("select id, email from users where email = ?;", email)

	var id int64
	var dbemail string
	for rows.Next() {
		err := rows.Scan(&id, &dbemail)
		if err != nil {
			t.Fatalf("DBからユーザを読み込めない")
		}
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

func TestFindOrCreateGithub(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	current_user, err := user.FindOrCreateGithub(token)
	if err != nil {
		t.Fatalf("Github経由で新規登録できない")
	}

	find_user, err := user.FindOrCreateGithub(token)
	if (find_user.Id != current_user.Id) || find_user.Id == 0 {
		t.Error("登録後にユーザを探せていない")
	}
}

func TestProjects(t *testing.T) {
	email := "project@example.com"
	password := "hogehoge"
	_ = user.Registration(email, password)
	mydb := &db.Database{}
	var database db.DB = mydb
	table := database.Init()
	rows, _ := table.Query("select id, email from users where email = ?;", email)

	var userid int64
	var dbemail string
	for rows.Next() {
		err := rows.Scan(&userid, &dbemail)
		if err != nil {
			t.Fatalf("DBからユーザを読み込めない")
		}
	}
	newProject := project.NewProject(0, userid, "project title")
	result := newProject.Save()
	if !result {
		t.Fatalf("プロジェクトが保存できない")
	}

	current_user := user.NewUser(userid, dbemail, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	projects := current_user.Projects()

	if projects[0].Id != newProject.Id {
		t.Error("ユーザとプロジェクトが関連づいていない")
	}
}

func TestCreateGithubUser(t *testing.T) {
	token := os.Getenv("TEST_TOKEN")
	newUser := user.NewUser(0, "", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
	result := newUser.CreateGithubUser(token)
	if !result {
		t.Fatalf("Github経由で新規登録できない")
	}

	mydb := &db.Database{}
	var database db.DB = mydb
	table := database.Init()
	rows, err := table.Query("select id, oauth_token from users where oauth_token = ?;", token)
	if err != nil {
		t.Fatalf("DB接続エラー")
	}
	var id int64
	var oauthToken sql.NullString
	for rows.Next() {
		err := rows.Scan(&id, &oauthToken)
		if err != nil {
			t.Fatalf("DBからユーザを読み込めない")
		}
	}
	if !oauthToken.Valid {
		t.Error("tokenに基づくユーザが登録されていない")
	}
	if id == int64(0) {
		t.Error("ユーザが保存されていない")
	}
}
