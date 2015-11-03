package project_test

import (
	"os"
	"database/sql"
	"../db"
	. "../project"
	"../user"
	"../list"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
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
		sql.Exec("truncate table lists;")
		sql.Close()
		os.Setenv("DB_NAME", currentdb)
	})

	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		uid, _ = user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		newProject = NewProject(0, uid, "title", "desc")
	})

	Describe("Save", func() {
		It("プロジェクトが登録できること", func() {
			result := newProject.Save()
			Expect(result).To(BeTrue())
			Expect(newProject.Id).NotTo(Equal(0))
		})
		It("ユーザとプロジェクトが関連付くこと", func() {
			_ = newProject.Save()
			rows, _ := table.Query("select id, user_id, title, description from projects where id = ?;", newProject.Id)

			var id int64
			var user_id sql.NullInt64
			var title string
			var description string

			for rows.Next() {
				err := rows.Scan(&id, &user_id, &title, &description)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(user_id.Valid).To(BeTrue())
			Expect(user_id.Int64).To(Equal(uid))
		})
	})

	Describe("Lists", func() {
		var (
			newList *list.ListStruct
			newProject *ProjectStruct
		)

		BeforeEach(func() {
			email := "lists@example.com"
			password := "hogehoge"
			user_id, _ := user.Registration(email, password)

			newProject = NewProject(0, user_id, "project title", "project desc")
			_ = newProject.Save()
			newList = list.NewList(0, newProject.Id, "list title", sql.NullString{})
			_ = newList.Save()
		})
		It("プロジェクトとリストが関連づいていること", func() {
			lists := newProject.Lists()
			Expect(lists).NotTo(BeEmpty())
			Expect(lists[0].Id).To(Equal(newList.Id))
		})

	})
})
