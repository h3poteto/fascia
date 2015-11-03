package task_test

import (
	"os"
	"database/sql"
	"../db"
	. "../task"
	"../project"
	"../list"
	"../user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		newList *list.ListStruct
		newTask *TaskStruct
		table *sql.DB
		currentdb string
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
		sql.Exec("truncate table tasks;")
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
		newProject := project.NewProject(0, uid, "title", "desc")
		newProject.Save()
		newList = list.NewList(0, newProject.Id, "list title", sql.NullString{})
		newList.Save()
		newTask = NewTask(0, newList.Id, "task title")
	})

	Describe("Save", func() {
		It("タスクが登録できること", func() {
			result := newTask.Save()
			Expect(result).To(BeTrue())
			Expect(newTask.Id).NotTo(Equal(0))
		})
		It("タスクとリストが関連づくこと", func() {
			_ = newTask.Save()
			rows, _ := table.Query("select id, list_id, title from tasks where id = ?;", newTask.Id)
			var id, list_id int64
			var title sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &list_id, &title)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(list_id).To(Equal(newTask.Id))
		})
	})
})
