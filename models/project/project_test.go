package project_test

import (
	"../../config"
	seed "../../db/seed"
	"../db"
	"../list"
	. "../project"
	"../repository"
	"../user"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	var (
		newProject *ProjectStruct
		uid        int64
		table      *sql.DB
	)

	BeforeEach(func() {
		seed.ListOptions()
	})
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table projects;")
		sql.Exec("truncate table repositories;")
		sql.Exec("truncate table lists;")
		sql.Close()
	})

	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		uid, _ = user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		newProject = NewProject(0, uid, "title", "desc", sql.NullInt64{}, true, true)
	})

	Describe("Create", func() {
		Context("when did not set repositoryID", func() {
			It("should create new project", func() {
				newProject, err := Create(uid, "new project", "description", 0, "", "", sql.NullString{})
				Expect(err).To(BeNil())
				Expect(len(newProject.Lists())).To(Equal(3))
				Expect(newProject.NoneList()).NotTo(BeNil())
				Expect(newProject.ShowIssues).To(BeTrue())
				Expect(newProject.ShowPullRequests).To(BeTrue())
			})
		})
	})

	Describe("Save", func() {
		It("should create project", func() {
			result := newProject.Save()
			Expect(result).To(BeTrue())
			Expect(newProject.Id).NotTo(Equal(0))
		})
		It("should relate user and project", func() {
			_ = newProject.Save()
			rows, _ := table.Query("select id, user_id, title, description from projects where id = ?;", newProject.Id)

			var id int64
			var user_id sql.NullInt64
			var title string
			var description string

			for rows.Next() {
				err := rows.Scan(&id, &user_id, &title, &description)
				if err != nil {
					panic(err)
				}
			}
			Expect(user_id.Valid).To(BeTrue())
			Expect(user_id.Int64).To(Equal(uid))
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			newProject.Save()
		})
		It("should set new value", func() {
			result := newProject.Update("newTitle", "newDescription", true, false)
			Expect(result).To(BeTrue())
			Expect(newProject.Title).To(Equal("newTitle"))
			Expect(newProject.Description).To(Equal("newDescription"))
			Expect(newProject.RepositoryId.Valid).To(BeFalse())
			Expect(newProject.ShowIssues).To(BeTrue())
			Expect(newProject.ShowPullRequests).To(BeFalse())
		})
	})

	Describe("Repository", func() {
		Context("when repository exist", func() {
			It("should relate project to repository", func() {
				repositoryId := int64(12345)
				newRepository := repository.NewRepository(0, repositoryId, "owner", "name")
				result := newRepository.Save()
				Expect(result).To(BeTrue())
				newProject.RepositoryId = sql.NullInt64{Int64: newRepository.Id, Valid: true}
				result = newProject.Save()

				Expect(result).To(BeTrue())
				Expect(newProject.Repository()).NotTo(BeNil())
				Expect(newProject.Repository().Id).To(Equal(newRepository.Id))
			})
		})
	})

	Describe("Lists", func() {
		var (
			newList, noneList *list.ListStruct
			newProject        *ProjectStruct
		)

		BeforeEach(func() {
			email := "lists@example.com"
			password := "hogehoge"
			user_id, _ := user.Registration(email, password)

			newProject = NewProject(0, user_id, "project title", "project desc", sql.NullInt64{}, true, true)
			_ = newProject.Save()
			newList = list.NewList(0, newProject.Id, newProject.UserId, "list title", "", sql.NullInt64{})
			_ = newList.Save(nil, nil)
			noneList = list.NewList(0, newProject.Id, newProject.UserId, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "", sql.NullInt64{})
			_ = noneList.Save(nil, nil)
		})
		It("should relate project and list", func() {
			lists := newProject.Lists()
			Expect(lists).NotTo(BeEmpty())
			Expect(lists[0].Id).To(Equal(newList.Id))
		})
		It("should not take none list", func() {
			lists := newProject.Lists()
			Expect(len(lists)).To(Equal(1))
		})

	})

	Describe("NoneList", func() {
		It("should contain only none list", func() {
			newProject, _ := Create(uid, "new project", "description", 0, "", "", sql.NullString{})
			noneList := newProject.NoneList()
			Expect(noneList.Title.String).To(Equal(config.Element("init_list").(map[interface{}]interface{})["none"].(string)))
		})
	})
})
