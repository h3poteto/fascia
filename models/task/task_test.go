package task_test

import (
	"../db"
	"../list"
	"../project"
	. "../task"
	"../user"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		newList    *list.ListStruct
		newTask    *TaskStruct
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
		newProject = project.NewProject(0, uid, "title", "desc", sql.NullInt64{}, true, true)
		newProject.Save()
		newList = list.NewList(0, newProject.Id, newProject.UserId, "list title", "", sql.NullInt64{})
		newList.Save(nil, nil)
		newTask = NewTask(0, newList.Id, newProject.Id, newList.UserId, sql.NullInt64{}, "task title", "task description", false, sql.NullString{})
	})

	Describe("Save", func() {
		It("can regist list", func() {
			result := newTask.Save(nil, nil)
			Expect(result).To(BeTrue())
			Expect(newTask.Id).NotTo(Equal(0))
		})
		It("should relate taks to list", func() {
			_ = newTask.Save(nil, nil)
			rows, _ := table.Query("select id, list_id, title from tasks where id = ?;", newTask.Id)
			var id, list_id int64
			var title sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &list_id, &title)
				if err != nil {
					panic(err)
				}
			}
			Expect(list_id).To(Equal(newTask.Id))
		})
		Context("when list do not have tasks", func() {
			It("should add display_index to task", func() {
				result := newTask.Save(nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.Id)
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
				existTask := NewTask(0, newList.Id, newProject.Id, newList.UserId, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save(nil, nil)
			})
			It("should set last display_index to task", func() {
				result := newTask.Save(nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.Id)
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
			list2 = list.NewList(0, newProject.Id, newProject.UserId, "list2", "", sql.NullInt64{})
			list2.Save(nil, nil)
		})
		Context("when destination list do not have tasks", func() {
			It("can move task", func() {
				result := newTask.ChangeList(list2.Id, nil, nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title from tasks where id = ?;", newTask.Id)
				var id, list_id int64
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title)
					if err != nil {
						panic(err)
					}
				}
				Expect(list_id).To(Equal(list2.Id))
			})
		})
		Context("when destination list have a task", func() {
			var (
				existTask *TaskStruct
			)
			JustBeforeEach(func() {
				existTask = NewTask(0, list2.Id, newProject.Id, list2.UserId, sql.NullInt64{}, "exist task title", "exist task description", false, sql.NullString{})
				existTask.Save(nil, nil)
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					result := newTask.ChangeList(list2.Id, nil, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
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
					result := newTask.ChangeList(list2.Id, &existTask.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
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
				existTask1 = NewTask(0, list2.Id, newProject.Id, list2.UserId, sql.NullInt64{}, "exist task title1", "exist task description1", false, sql.NullString{})
				existTask1.Save(nil, nil)
				existTask2 = NewTask(0, list2.Id, newProject.Id, list2.UserId, sql.NullInt64{}, "exist task title2", "exist task description2", false, sql.NullString{})
				existTask2.Save(nil, nil)
			})
			Context("when send nil", func() {
				It("should add task to end of list", func() {
					result := newTask.ChangeList(list2.Id, nil, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
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
					result := newTask.ChangeList(list2.Id, &existTask1.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
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
					result := newTask.ChangeList(list2.Id, &existTask2.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
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
					result := newTask.ChangeList(list2.Id, &existTask2.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, existTask2.Id)
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
})
