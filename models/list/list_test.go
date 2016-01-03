package list_test

import (
	"../db"
	. "../list"
	"../project"
	"../task"
	"../user"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		newList    *ListStruct
		newProject *project.ProjectStruct
		table      *sql.DB
	)
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table projects;")
		sql.Exec("truncate table lists;")
		sql.Exec("truncate table tasks;")
		sql.Close()
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
		newList = NewList(0, newProject.Id, newProject.UserId, "list title", "")
	})

	Describe("Save", func() {
		It("can registrate list", func() {
			result := newList.Save(nil, nil)
			Expect(result).To(BeTrue())
			Expect(newList.Id).NotTo(Equal(0))
		})
		It("should relate list to project", func() {
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
		It("should find list which related project", func() {
			newList.Save(nil, nil)
			findList := FindList(newProject.Id, newList.Id)
			Expect(findList).To(Equal(newList))
		})
	})

	Describe("Tasks", func() {
		var newTask *task.TaskStruct
		JustBeforeEach(func() {
			newList.Save(nil, nil)
			newTask = task.NewTask(0, newList.Id, newList.UserId, sql.NullInt64{}, "task", "description")
			newTask.Save(nil, nil)
		})
		It("should related task to list", func() {
			tasks := newList.Tasks()
			Expect(tasks).NotTo(BeEmpty())
			Expect(tasks[0].Id).To(Equal(newTask.Id))
		})

	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			newList.Save(nil, nil)
		})
		It("should update list", func() {
			newTitle := "newTitle"
			newColor := "newColor"
			newList.Update(nil, nil, &newTitle, &newColor)
			findList := FindList(newList.ProjectId, newList.Id)
			Expect(findList.Title.String).To(Equal(newTitle))
			Expect(findList.Color.String).To(Equal(newColor))
		})
	})
})
