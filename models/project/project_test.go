package project_test

import (
	"os"
	"database/sql"
	"../db"
	. "../project"
	"../user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectSave", func() {
	var (
		newProject *ProjectStruct
		currentdb string
		uid int64
		table *sql.DB
	)

	BeforeEach(func() {
		testdb := os.Getenv("DB_TEST_NAME")
		currentdb = os.Getenv("DB_NAME")
		os.Setenv("DB_NAME", testdb)
	})
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table projects;")
		sql.Close()
		os.Setenv("DB_NAME", currentdb)
	})

	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		_ = user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		rows, _ := table.Query("select id from users where email = ?;", email)
		for rows.Next() {
			err := rows.Scan(&uid)
			if err != nil {
				panic(err.Error())
			}
		}
		newProject = NewProject(0, uid, "title")
	})

	Describe("Save", func() {
		It("プロジェクトが登録できること", func() {
			result := newProject.Save()
			Expect(result).To(BeTrue())
			Expect(newProject.Id).NotTo(Equal(0))
		})
		It("ユーザとプロジェクトが関連付くこと", func() {
			_ = newProject.Save()
			rows, _ := table.Query("select id, user_id, title from projects where id = ?;", newProject.Id)

			var id int64
			var user_id sql.NullInt64
			var title string

			for rows.Next() {
				err := rows.Scan(&id, &user_id, &title)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(user_id.Valid).To(BeTrue())
			Expect(user_id.Int64).To(Equal(uid))
		})
	})
})
