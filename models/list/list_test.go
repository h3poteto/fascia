package list_test

import (
	"os"
	"database/sql"
	"../db"
	. "../list"
	"../project"
	"../user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListSave", func() {
	var (
		newList *ListStruct
		newProject *project.ProjectStruct
		currentdb string
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
		sql.Exec("truncate table lists;")
		sql.Close()
		os.Setenv("DB_NAME", currentdb)
	})

	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		uid, _ := user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		newProject = project.NewProject(0, uid, "title")
		newProject.Save()
		newList = NewList(0, newProject.Id, "list title")
	})

	Describe("Save", func() {
		It("リストが登録できること", func() {
			result := newList.Save()
			Expect(result).To(BeTrue())
			Expect(newList.Id).NotTo(Equal(0))
		})
		It("プロジェクトとリストが関連づくこと", func() {
			_ = newList.Save()
			rows, _ := table.Query("select id, project_id, title from lists where id = ?;", newList.Id)
			var id int64
			var project_id int64
			var title sql.NullString

			for rows.Next() {
				err := rows.Scan(&id, &project_id, &title)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(project_id).To(Equal(newProject.Id))
		})
	})
})
