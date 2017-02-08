package list_test

import (
	"database/sql"

	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/entities/list"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var (
		newList        *List
		projectService *services.Project
		database       *sql.DB
	)

	BeforeEach(func() {
		seed.Seeds()
		email := "save@example.com"
		password := "hogehoge"
		user, _ := handlers.RegistrationUser(email, password, password)
		database = db.SharedInstance().Connection
		projectService, _ = handlers.CreateProject(user.UserEntity.UserModel.ID, "title", "desc", 0, sql.NullString{})
		newList = New(0, projectService.ProjectEntity.ProjectModel.ID, projectService.ProjectEntity.ProjectModel.UserID, "list title", "", sql.NullInt64{}, false)
	})

	Describe("Save", func() {
		It("can registrate list", func() {
			err := newList.Save(nil)
			Expect(err).To(BeNil())
			Expect(newList.ListModel.ID).NotTo(Equal(0))
		})
		It("should relate list to project", func() {
			_ = newList.Save(nil)
			rows, _ := database.Query("select id, project_id, title from lists where id = ?;", newList.ListModel.ID)
			var id int64
			var projectID int64
			var title sql.NullString

			for rows.Next() {
				err := rows.Scan(&id, &projectID, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(projectID).To(Equal(projectService.ProjectEntity.ProjectModel.ID))
		})
	})

	Describe("FindList", func() {
		It("should find list which related project", func() {
			newList.Save(nil)
			findList, err := FindByID(projectService.ProjectEntity.ProjectModel.ID, newList.ListModel.ID)
			Expect(err).To(BeNil())
			Expect(findList.ListModel).To(Equal(newList.ListModel))
		})
	})

	Describe("Tasks", func() {
		var taskService *services.Task
		JustBeforeEach(func() {
			newList.Save(nil)
			taskService = services.NewTask(0, newList.ListModel.ID, projectService.ProjectEntity.ProjectModel.ID, newList.ListModel.UserID, sql.NullInt64{}, "task", "description", false, sql.NullString{})
			taskService.Save()
		})
		It("should related task to list", func() {
			tasks, err := newList.Tasks()
			Expect(err).To(BeNil())
			Expect(tasks).NotTo(BeEmpty())
			Expect(tasks[0].TaskModel.ID).To(Equal(taskService.TaskEntity.TaskModel.ID))
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
				findList, err := FindByID(newList.ListModel.ProjectID, newList.ListModel.ID)
				Expect(err).To(BeNil())
				Expect(findList.ListModel.Title.String).To(Equal(newTitle))
				Expect(findList.ListModel.Color.String).To(Equal(newColor))
			})
		})
		Context("have list_option", func() {
			It("should update list and have list_option", func() {
				newTitle := "newTitle"
				newColor := "newColor"
				listOption, _ := services.FindListOptionByAction("close")
				newList.Update(newTitle, newColor, listOption.ListOptionEntity.ListOptionModel.ID)
				findList, err := FindByID(newList.ListModel.ProjectID, newList.ListModel.ID)
				Expect(err).To(BeNil())
				Expect(findList.ListModel.Title.String).To(Equal(newTitle))
				Expect(findList.ListModel.Color.String).To(Equal(newColor))
				Expect(findList.ListModel.ListOptionID.Int64).To(Equal(listOption.ListOptionEntity.ListOptionModel.ID))
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
			Expect(newList.ListModel.IsHidden).To(BeTrue())
			l, _ := FindByID(projectService.ProjectEntity.ProjectModel.ID, newList.ListModel.ID)
			Expect(l.ListModel.IsHidden).To(BeTrue())
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
			Expect(newList.ListModel.IsHidden).To(BeFalse())
			l, _ := FindByID(projectService.ProjectEntity.ProjectModel.ID, newList.ListModel.ID)
			Expect(l.ListModel.IsHidden).To(BeFalse())
		})
	})

	Describe("DeleteTasks", func() {
		var taskService *services.Task
		JustBeforeEach(func() {
			newList.Save(nil)
			taskService = services.NewTask(0, newList.ListModel.ID, projectService.ProjectEntity.ProjectModel.ID, newList.ListModel.UserID, sql.NullInt64{}, "task", "description", false, sql.NullString{})
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
			Expect(newList.ListModel).To(BeNil())
		})
	})
})
