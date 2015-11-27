package list_test

import (
	"../db"
	. "../list"
	"../project"
	"../task"
	"../user"
	"database/sql"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		newList    *ListStruct
		newProject *project.ProjectStruct
		currentdb  string
		table      *sql.DB
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
		newProject = project.NewProject(0, uid, "title", "desc")
		newProject.Save()
		newList = NewList(0, newProject.Id, newProject.UserId.Int64, "list title", "")
	})

	Describe("Save", func() {
		It("リストが登録できること", func() {
			result := newList.Save(nil, nil)
			Expect(result).To(BeTrue())
			Expect(newList.Id).NotTo(Equal(0))
		})
		It("プロジェクトとリストが関連づくこと", func() {
			_ = newList.Save(nil, nil)
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

	Describe("FindList", func() {
		It("プロジェクトに関連づいたリストが見つかること", func() {
			newList.Save(nil, nil)
			findList := FindList(newProject.Id, newList.Id)
			Expect(findList).To(Equal(newList))
		})
	})

	Describe("Tasks", func() {
		var newTask *task.TaskStruct
		JustBeforeEach(func() {
			newList.Save(nil, nil)
			newTask = task.NewTask(0, newList.Id, newList.UserId, sql.NullInt64{}, "task")
			newTask.Save(nil, nil)
		})
		It("taskが関連づくこと", func() {
			tasks := newList.Tasks()
			Expect(tasks).NotTo(BeEmpty())
			Expect(tasks[0].Id).To(Equal(newTask.Id))
		})

	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			newList.Save(nil, nil)
		})
		It("リストが更新できること", func() {
			newTitle := "newTitle"
			newColor := "newColor"
			newList.Update(nil, nil, &newTitle, &newColor)
			findList := FindList(newList.ProjectId, newList.Id)
			Expect(findList.Title.String).To(Equal(newTitle))
			Expect(findList.Color.String).To(Equal(newColor))
		})
	})
})
