package project

import (
	"os"
	"testing"
	"database/sql"
	"../models/db"
	"../models/project"
	"../models/user"
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

func TestSave(t *testing.T) {
	email := "save@example.com"
	password := "hogehoge"

	_ = user.Registration(email, password)

	mydb := &db.Database{}
	var database db.DB = mydb
	table := database.Init()
	rows, _ := table.Query("select id from users where email = ?;", email)
	var uid int64
	for rows.Next() {
		err := rows.Scan(&uid)
		if err != nil {
			t.Fatalf("DBからユーザを読み込めない")
		}
	}

	newProject := project.NewProject(0, uid, "title")
	result := newProject.Save()
	if !result {
		t.Error("プロジェクトが登録できない")
	}

	if newProject.Id == 0 {
		t.Error("プロジェクトが登録できない")
	}

	rows, _ = table.Query("select id, user_id, title from projects where id = ?;", newProject.Id)

	var id int64
	var user_id sql.NullInt64
	var title string

	for rows.Next() {
		err := rows.Scan(&id, &user_id, &title)
		if err != nil {
			t.Fatalf("DBからユーザを読み込めない")
		}
	}

	if !user_id.Valid {
		t.Error("ユーザIDが入っていない")
	}

	if user_id.Int64 != uid {
		t.Error("プロジェクトとユーザが関連づいていない")
	}

}
