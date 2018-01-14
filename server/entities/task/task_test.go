package task_test

import (
	"database/sql"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/commands/board"
	. "github.com/h3poteto/fascia/server/entities/task"
	"github.com/h3poteto/fascia/server/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		listService    *board.List
		newTask        *Task
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
		listService = handlers.NewList(0, projectService.ProjectEntity.ID, projectService.ProjectEntity.UserID, "list title", "", sql.NullInt64{}, false)
		listService.Save()
		newTask = New(0, listService.ListEntity.ID, projectService.ProjectEntity.ID, listService.ListEntity.UserID, sql.NullInt64{}, "task title", "task description", false, sql.NullString{})
	})

	Describe("Save", func() {
		It("can regist list", func() {
			err := newTask.Save()
			Expect(err).To(BeNil())
			Expect(newTask.ID).NotTo(Equal(0))
		})
		It("should relate taks to list", func() {
			_ = newTask.Save()
			rows, _ := db.Query("select id, list_id, title from tasks where id = ?;", newTask.ID)
			var id, listID int64
			var title sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &listID, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(listID).To(Equal(listService.ListEntity.ID))
		})
		Context("when list do not have tasks", func() {
			It("should add display_index to task", func() {
				err := newTask.Save()
				Expect(err).To(BeNil())
				rows, _ := db.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.ID)
				var id, listID int64
				var displayIndex int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &listID, &title, &displayIndex)
					if err != nil {
						panic(err)
					}
				}
				Expect(displayIndex).To(Equal(1))
			})
		})
		Context("when list have tasks", func() {
			JustBeforeEach(func() {
				existTask := New(0, listService.ListEntity.ID, projectService.ProjectEntity.ID, listService.ListEntity.UserID, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save()
			})
			It("should set last display_index to task", func() {
				err := newTask.Save()
				Expect(err).To(BeNil())
				rows, _ := db.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.ID)
				var id, listID int64
				var displayIndex int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &listID, &title, &displayIndex)
					if err != nil {
						panic(err)
					}
				}
				Expect(displayIndex).To(Equal(2))
			})
		})
	})

	Describe("ChangeList", func() {
		var (
			secondaryList *board.List
		)
		BeforeEach(func() {
			newTask.Save()
			secondaryList = handlers.NewList(0, projectService.ProjectEntity.ID, projectService.ProjectEntity.UserID, "list2", "", sql.NullInt64{}, false)
			secondaryList.Save()
		})
		Context("when destination list do not have tasks", func() {
			It("can move task", func() {
				isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, nil)
				Expect(err).To(BeNil())
				Expect(isReorder).To(BeFalse())
				rows, _ := db.Query("select id, list_id, title from tasks where id = ?;", newTask.ID)
				var id, listID int64
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &listID, &title)
					if err != nil {
						panic(err)
					}
				}
				Expect(listID).To(Equal(secondaryList.ListEntity.ID))
			})
		})
		Context("when destination list have a task", func() {
			var (
				existTask *Task
			)
			BeforeEach(func() {
				existTask = New(0, secondaryList.ListEntity.ID, projectService.ProjectEntity.ID, secondaryList.ListEntity.UserID, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save()
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, nil)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, newTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(2))
				})
			})
			Context("when add task before exist task", func() {
				It("should add task to top of list", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, &existTask.ID)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, newTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(1))
				})
			})
		})
		Context("when destination list have tasks", func() {
			var firstExistTask, secondExistTask *Task
			BeforeEach(func() {
				firstExistTask = New(0, secondaryList.ListEntity.ID, projectService.ProjectEntity.ID, secondaryList.ListEntity.UserID, sql.NullInt64{}, "exist task title1", "exist task description1", false, sql.NullString{})
				firstExistTask.Save()
				secondExistTask = New(0, secondaryList.ListEntity.ID, projectService.ProjectEntity.ID, secondaryList.ListEntity.UserID, sql.NullInt64{}, "exist task title2", "exist task description2", false, sql.NullString{})
				secondExistTask.Save()
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, nil)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, newTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(3))
				})
			})
			Context("when send task to top of list", func() {
				It("should add task to top of list", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, &firstExistTask.ID)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, newTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(1))
				})
			})
			Context("when send task to mid-flow", func() {
				It("should add task to mid-flow", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, &secondExistTask.ID)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, newTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(2))
				})
				It("other tasks should be pushed out", func() {
					isReorder, err := newTask.ChangeList(secondaryList.ListEntity.ID, &secondExistTask.ID)
					Expect(err).To(BeNil())
					Expect(isReorder).To(BeFalse())
					rows, _ := db.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", secondaryList.ListEntity.ID, secondExistTask.ID)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err)
						}
					}
					Expect(displayIndex).To(Equal(3))
				})
			})
		})
	})

	Describe("Delete", func() {
		Context("when a task does not relate issue", func() {
			It("can delete task", func() {
				newTask.Save()
				err := newTask.Delete()
				Expect(err).To(BeNil())
			})
		})
		Context("when a task relate issue", func() {
			It("cannot delete task", func() {
				newIssueTask := New(0, listService.ListEntity.ID, projectService.ProjectEntity.ID, listService.ListEntity.UserID, sql.NullInt64{Int64: 1, Valid: true}, "issue title", "issue description", false, sql.NullString{})
				newIssueTask.Save()
				err := newIssueTask.Delete()
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
