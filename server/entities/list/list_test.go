package list_test

import (
	"database/sql"

	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/commands/board"
	. "github.com/h3poteto/fascia/server/entities/list"
	"github.com/h3poteto/fascia/server/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		newList        *List
		projectService *board.Project
		db             *sql.DB
	)

	BeforeEach(func() {
		seed.Seeds()
		email := "save@example.com"
		password := "hogehoge"
		user, _ := handlers.RegistrationUser(email, password, password)
		db = database.SharedInstance().Connection
		projectService, _ = handlers.CreateProject(user.UserEntity.ID, "title", "desc", 0, sql.NullString{})
		newList = New(0, projectService.ProjectEntity.ID, projectService.ProjectEntity.UserID, "list title", "", sql.NullInt64{}, false)
	})

	Describe("Save", func() {
		It("can registrate list", func() {
			err := newList.Save(nil)
			Expect(err).To(BeNil())
			Expect(newList.ID).NotTo(Equal(0))
		})
		It("should relate list to project", func() {
			_ = newList.Save(nil)
			rows, _ := db.Query("select id, project_id, title from lists where id = ?;", newList.ID)
			var id int64
			var projectID int64
			var title sql.NullString

			for rows.Next() {
				err := rows.Scan(&id, &projectID, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(projectID).To(Equal(projectService.ProjectEntity.ID))
		})
	})

	Describe("FindList", func() {
		It("should find list which related project", func() {
			newList.Save(nil)
			findList, err := FindByID(projectService.ProjectEntity.ID, newList.ID)
			Expect(err).To(BeNil())
			Expect(findList.ID).To(Equal(newList.ID))
		})
	})

	Describe("Tasks", func() {
		var taskService *board.Task
		JustBeforeEach(func() {
			newList.Save(nil)
			taskService = board.NewTask(0, newList.ID, projectService.ProjectEntity.ID, newList.UserID, sql.NullInt64{}, "task", "description", false, sql.NullString{})
			taskService.Save()
		})
		It("should related task to list", func() {
			tasks, err := newList.Tasks()
			Expect(err).To(BeNil())
			Expect(tasks).NotTo(BeEmpty())
			Expect(tasks[0].ID).To(Equal(taskService.TaskEntity.ID))
		})

	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			newList.Save(nil)
		})
		Context("does not have list_option", func() {
			It("should update list", func() {
				newTitle := "newTitle"
				newColor := "newColor"
				optionID := int64(0)
				newList.Update(newTitle, newColor, optionID)
				findList, err := FindByID(newList.ProjectID, newList.ID)
				Expect(err).To(BeNil())
				Expect(findList.Title.String).To(Equal(newTitle))
				Expect(findList.Color.String).To(Equal(newColor))
			})
		})
		Context("have list_option", func() {
			It("should update list and have list_option", func() {
				newTitle := "newTitle"
				newColor := "newColor"
				listOption, _ := board.FindListOptionByAction("close")
				newList.Update(newTitle, newColor, listOption.ListOptionEntity.ID)
				findList, err := FindByID(newList.ProjectID, newList.ID)
				Expect(err).To(BeNil())
				Expect(findList.Title.String).To(Equal(newTitle))
				Expect(findList.Color.String).To(Equal(newColor))
				Expect(findList.ListOptionID.Int64).To(Equal(listOption.ListOptionEntity.ID))
			})
		})
	})

	Describe("Hide", func() {
		JustBeforeEach(func() {
			newList.Save(nil)
		})
		It("should hidden list", func() {
			err := newList.Hide()
			Expect(err).To(BeNil())
			Expect(newList.IsHidden).To(BeTrue())
			l, _ := FindByID(projectService.ProjectEntity.ID, newList.ID)
			Expect(l.IsHidden).To(BeTrue())
		})
	})

	Describe("Display", func() {
		JustBeforeEach(func() {
			newList.Save(nil)
			newList.Hide()
		})
		It("should display list", func() {
			err := newList.Display()
			Expect(err).To(BeNil())
			Expect(newList.IsHidden).To(BeFalse())
			l, _ := FindByID(projectService.ProjectEntity.ID, newList.ID)
			Expect(l.IsHidden).To(BeFalse())
		})
	})

	Describe("DeleteTasks", func() {
		var taskService *board.Task
		JustBeforeEach(func() {
			newList.Save(nil)
			taskService = board.NewTask(0, newList.ID, projectService.ProjectEntity.ID, newList.UserID, sql.NullInt64{}, "task", "description", false, sql.NullString{})
			taskService.Save()
		})
		It("should delete all tasks", func() {
			err := newList.DeleteTasks()
			Expect(err).To(BeNil())
			tasks, _ := newList.Tasks()
			Expect(len(tasks)).To(Equal(0))
		})
	})

	Describe("Delete", func() {
		JustBeforeEach(func() {
			newList.Save(nil)
		})
		It("should delete list", func() {
			err := newList.Delete()
			Expect(err).To(BeNil())
		})
	})
})
