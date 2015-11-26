package task_test

import (
	"../db"
	"../list"
	"../project"
	. "../task"
	"../user"
	"database/sql"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		newList    *list.ListStruct
		newTask    *TaskStruct
		newProject *project.ProjectStruct
		table      *sql.DB
		currentdb  string
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
		newList = list.NewList(0, newProject.Id, newProject.UserId.Int64, "list title", "")
		newList.Save(nil, nil)
		newTask = NewTask(0, newList.Id, newList.UserId, sql.NullInt64{}, "task title")
	})

	Describe("Save", func() {
		It("タスクが登録できること", func() {
			result := newTask.Save(nil, nil)
			Expect(result).To(BeTrue())
			Expect(newTask.Id).NotTo(Equal(0))
		})
		It("タスクとリストが関連づくこと", func() {
			_ = newTask.Save(nil, nil)
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
		Context("リストにタスクがないとき", func() {
			It("display_indexが追加されていること", func() {
				result := newTask.Save(nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.Id)
				var id, list_id int64
				var display_index int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title, &display_index)
					if err != nil {
						panic(err.Error())
					}
				}
				Expect(display_index).To(Equal(1))
			})
		})
		Context("リストにタスクがあるとき", func() {
			JustBeforeEach(func() {
				existTask := NewTask(0, newList.Id, newList.UserId, sql.NullInt64{}, "exist task title")
				existTask.Save(nil, nil)
			})
			It("display_indexがラストになっていること", func() {
				result := newTask.Save(nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title, display_index from tasks where id = ?;", newTask.Id)
				var id, list_id int64
				var display_index int
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title, &display_index)
					if err != nil {
						panic(err.Error())
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
			list2 = list.NewList(0, newProject.Id, newProject.UserId.Int64, "list2", "")
			list2.Save(nil, nil)
		})
		Context("移動先リストにタスクがないとき", func() {
			It("タスクが移動できること", func() {
				result := newTask.ChangeList(list2.Id, nil, nil, nil)
				Expect(result).To(BeTrue())
				rows, _ := table.Query("select id, list_id, title from tasks where id = ?;", newTask.Id)
				var id, list_id int64
				var title sql.NullString
				for rows.Next() {
					err := rows.Scan(&id, &list_id, &title)
					if err != nil {
						panic(err.Error())
					}
				}
				Expect(list_id).To(Equal(list2.Id))
			})
		})
		Context("移動先リストにタスクがひとつだけの時", func() {
			var (
				existTask *TaskStruct
			)
			JustBeforeEach(func() {
				existTask = NewTask(0, list2.Id, list2.UserId, sql.NullInt64{}, "exist task title")
				existTask.Save(nil, nil)
			})
			Context("nil順位を渡した時", func() {
				It("末尾に挿入されること", func() {
					result := newTask.ChangeList(list2.Id, nil, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(2))
				})
			})
			Context("存在するタスクの前に入れたいとき", func() {
				It("先頭に挿入されること", func() {
					result := newTask.ChangeList(list2.Id, &existTask.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(1))
				})
			})
		})
		Context("移動先リストに複数タスクがあるとき", func() {
			var existTask1, existTask2 *TaskStruct
			JustBeforeEach(func() {
				existTask1 = NewTask(0, list2.Id, list2.UserId, sql.NullInt64{}, "exist task title1")
				existTask1.Save(nil, nil)
				existTask2 = NewTask(0, list2.Id, list2.UserId, sql.NullInt64{}, "exist task title2")
				existTask2.Save(nil, nil)
			})
			Context("nil順位を渡した時", func() {
				It("末尾に挿入されること", func() {
					result := newTask.ChangeList(list2.Id, nil, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(3))
				})
			})
			Context("先頭に入れたいとき", func() {
				It("先頭に挿入されること", func() {
					result := newTask.ChangeList(list2.Id, &existTask1.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(1))
				})
			})
			Context("途中に入れたいとき", func() {
				It("途中に挿入されること", func() {
					result := newTask.ChangeList(list2.Id, &existTask2.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, newTask.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(2))
				})
				It("他のタスクが押し出されていること", func() {
					result := newTask.ChangeList(list2.Id, &existTask2.Id, nil, nil)
					Expect(result).To(BeTrue())
					rows, _ := table.Query("select id, title, display_index from tasks where list_id = ? and id = ?;", list2.Id, existTask2.Id)
					var id int64
					var displayIndex int
					var title sql.NullString
					for rows.Next() {
						err := rows.Scan(&id, &title, &displayIndex)
						if err != nil {
							panic(err.Error())
						}
					}
					Expect(displayIndex).To(Equal(3))
				})
			})
		})
	})
})
