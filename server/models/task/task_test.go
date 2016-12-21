package task_test

import (
	"database/sql"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/list"
	"github.com/h3poteto/fascia/models/project"
	. "github.com/h3poteto/fascia/models/task"
	"github.com/h3poteto/fascia/models/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		newList    *list.ListStruct
		newTask    *TaskStruct
		newProject *project.ProjectStruct
		database   *sql.DB
	)
	BeforeEach(func() {
		seed.Seeds()
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
		database.Exec("truncate table lists;")
		database.Exec("truncate table tasks;")
		database.Exec("truncate table list_options;")
	})
	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		uid, _ := user.Registration(email, password, password)
		database = db.SharedInstance().Connection
		newProject, _ = project.Create(uid, "title", "desc", 0, sql.NullString{})
		newList = list.NewList(0, newProject.ID, newProject.UserID, "list title", "", sql.NullInt64{}, false)
		newList.Save(nil, nil)
		newTask = NewTask(0, newList.ID, newProject.ID, newList.UserID, sql.NullInt64{}, "task title", "task description", false, sql.NullString{})
	})

	Describe("Save", func() {
		It("can regist list", func() {
			err := newTask.Save(nil, nil)
			Expect(err).To(BeNil())
			Expect(newTask.ID).NotTo(Equal(0))
		})
		It("should relate taks to list", func() {
			_ = newTask.Save(nil, nil)
			rows, _ := database.Query("select id, list_id, title from tasks where id = ?;", newTask.ID)
			var id, list_id int64
			var title sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &list_id, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(list_id).To(Equal(newList.ID))
		})
		Context("when list do not have tasks", func() {
			It("should add display_index to task", func() {
				err := newTask.Save(nil, nil)
				Expect(err).To(BeNil())
				rows, _ := database.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.ID)
				var id, list_id int64
				var display_index int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title, &display_index)
					if err != nil {
						panic(err)
					}
				}
				Expect(display_index).To(Equal(1))
			})
		})
		Context("when list have tasks", func() {
			JustBeforeEach(func() {
				existTask := NewTask(0, newList.ID, newProject.ID, newList.UserID, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save(nil, nil)
			})
			It("should set last display_index to task", func() {
				err := newTask.Save(nil, nil)
				Expect(err).To(BeNil())
				rows, _ := database.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.ID)
				var id, list_id int64
				var display_index int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title, &display_index)
					if err != nil {
						panic(err)
					}
				}
				Expect(display_index).To(Equal(2))
			})
		})
	})

	Describe("ChangeList", func() {
		var (
			list2 *list.ListStruct
		)
		JustBeforeEach(func() {
			newTask.Save(nil, nil)
			list2 = list.NewList(0, newProject.ID, newProject.UserID, "list2", "", sql.NullInt64{}, false)
			list2.Save(nil, nil)
		})
		Context("when destination list do not have tasks", func() {
			It("can move task", func() {
				err := newTask.ChangeList(list2.ID, nil, nil, nil)
				Expect(err).To(BeNil())
				rows, _ := database.Query("select id, list_id, title from tasks where id = ?;", newTask.ID)
				var id, list_id int64
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title)
					if err != nil {
						panic(err)
					}
				}
				Expect(list_id).To(Equal(list2.ID))
			})
		})
		Context("when destination list have a task", func() {
			var (
				existTask *TaskStruct
			)
			JustBeforeEach(func() {
				existTask = NewTask(0, list2.ID, newProject.ID, list2.UserID, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save(nil, nil)
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					err := newTask.ChangeList(list2.ID, nil, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, newTask.ID)
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
					err := newTask.ChangeList(list2.ID, &existTask.ID, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, newTask.ID)
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
			var existTask1, existTask2 *TaskStruct
			JustBeforeEach(func() {
				existTask1 = NewTask(0, list2.ID, newProject.ID, list2.UserID, sql.NullInt64{}, "exist task title1", "exist task description1", false, sql.NullString{})
				existTask1.Save(nil, nil)
				existTask2 = NewTask(0, list2.ID, newProject.ID, list2.UserID, sql.NullInt64{}, "exist task title2", "exist task description2", false, sql.NullString{})
				existTask2.Save(nil, nil)
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					err := newTask.ChangeList(list2.ID, nil, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, newTask.ID)
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
					err := newTask.ChangeList(list2.ID, &existTask1.ID, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, newTask.ID)
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
					err := newTask.ChangeList(list2.ID, &existTask2.ID, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, newTask.ID)
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
					err := newTask.ChangeList(list2.ID, &existTask2.ID, nil, nil)
					Expect(err).To(BeNil())
					rows, _ := database.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.ID, existTask2.ID)
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
				newTask.Save(nil, nil)
				err := newTask.Delete()
				Expect(err).To(BeNil())
				Expect(newTask.ID).To(BeEquivalentTo(int64(0)))
			})
		})
		Context("when a task relate issue", func() {
			It("cannot delete task", func() {
				newIssueTask := NewTask(0, newList.ID, newProject.ID, newList.UserID, sql.NullInt64{Int64: 1, Valid: true}, "issue title", "issue description", false, sql.NullString{})
				newIssueTask.Save(nil, nil)
				err := newIssueTask.Delete()
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
