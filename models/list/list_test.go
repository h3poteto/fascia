package list_test

import (
	"database/sql"

	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	. "github.com/h3poteto/fascia/models/list"
	"github.com/h3poteto/fascia/models/list_option"
	"github.com/h3poteto/fascia/models/project"
	"github.com/h3poteto/fascia/models/task"
	"github.com/h3poteto/fascia/models/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		newList    *ListStruct
		newProject *project.ProjectStruct
		database   *sql.DB
	)
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
		database.Exec("truncate table lists;")
		database.Exec("truncate table tasks;")
		database.Exec("truncate table list_options;")
	})

	JustBeforeEach(func() {
		seed.Seeds()
		email := "save@example.com"
		password := "hogehoge"
		uid, _ := user.Registration(email, password, password)
		database = db.SharedInstance().Connection
		newProject, _ = project.Create(uid, "title", "desc", 0, sql.NullString{})
		newList = NewList(0, newProject.ID, newProject.UserID, "list title", "", sql.NullInt64{}, false)
	})

	Describe("Save", func() {
		It("can registrate list", func() {
			err := newList.Save(nil, nil)
			Expect(err).To(BeNil())
			Expect(newList.ID).NotTo(Equal(0))
		})
		It("should relate list to project", func() {
			_ = newList.Save(nil, nil)
			rows, _ := database.Query("select id, project_id, title from lists where id = ?;", newList.ID)
			var id int64
			var project_id int64
			var title sql.NullString

			for rows.Next() {
				err := rows.Scan(&id, &project_id, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(project_id).To(Equal(newProject.ID))
		})
	})

	Describe("FindList", func() {
		It("should find list which related project", func() {
			newList.Save(nil, nil)
			findList, err := FindList(newProject.ID, newList.ID)
			Expect(err).To(BeNil())
			Expect(findList).To(Equal(newList))
		})
	})

	Describe("Tasks", func() {
		var newTask *task.TaskStruct
		JustBeforeEach(func() {
			newList.Save(nil, nil)
			newTask = task.NewTask(0, newList.ID, newProject.ID, newList.UserID, sql.NullInt64{}, "task", "description", false, sql.NullString{})
			newTask.Save(nil, nil)
		})
		It("should related task to list", func() {
			tasks, err := newList.Tasks()
			Expect(err).To(BeNil())
			Expect(tasks).NotTo(BeEmpty())
			Expect(tasks[0].ID).To(Equal(newTask.ID))
		})

	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			newList.Save(nil, nil)
		})
		Context("does not have list_option", func() {
			It("should update list", func() {
				newTitle := "newTitle"
				newColor := "newColor"
				optionID := int64(0)
				newList.Update(nil, nil, &newTitle, &newColor, &optionID)
				findList, err := FindList(newList.ProjectID, newList.ID)
				Expect(err).To(BeNil())
				Expect(findList.Title.String).To(Equal(newTitle))
				Expect(findList.Color.String).To(Equal(newColor))
			})
		})
		Context("have list_option", func() {
			It("should update list and have list_option", func() {
				newTitle := "newTitle"
				newColor := "newColor"
				listOption, _ := list_option.FindByAction("close")
				newList.Update(nil, nil, &newTitle, &newColor, &listOption.ID)
				findList, err := FindList(newList.ProjectID, newList.ID)
				Expect(err).To(BeNil())
				Expect(findList.Title.String).To(Equal(newTitle))
				Expect(findList.Color.String).To(Equal(newColor))
				Expect(findList.ListOptionID.Int64).To(Equal(listOption.ID))
			})
		})
	})

	Describe("Hide", func() {
		JustBeforeEach(func() {
			newList.Save(nil, nil)
		})
		It("should hidden list", func() {
			err := newList.Hide()
			Expect(err).To(BeNil())
			Expect(newList.IsHidden).To(BeTrue())
			l, _ := FindList(newProject.ID, newList.ID)
			Expect(l.IsHidden).To(BeTrue())
		})
	})

	Describe("Display", func() {
		JustBeforeEach(func() {
			newList.Save(nil, nil)
			newList.Hide()
		})
		It("should display list", func() {
			err := newList.Display()
			Expect(err).To(BeNil())
			Expect(newList.IsHidden).To(BeFalse())
			l, _ := FindList(newProject.ID, newList.ID)
			Expect(l.IsHidden).To(BeFalse())
		})
	})

})
