package project_test

import (
	"../../config"
	seed "../../db/seed"
	"../db"
	"../list"
	. "../project"
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
			It("should relate user and project", func() {
				newProject, _ = Create(uid, "new project", "description", 0, "", "", sql.NullString{})
				rows, _ := table.Query("select id, user_id, title, description from projects where id = ?;", newProject.ID)

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
	})

	Describe("Update", func() {
		BeforeEach(func() {
			newProject, _ = Create(uid, "new project", "description", 0, "", "", sql.NullString{})
		})
		It("should set new value", func() {
			err := newProject.Update("newTitle", "newDescription", true, false)
			Expect(err).To(BeNil())
			Expect(newProject.Title).To(Equal("newTitle"))
			Expect(newProject.Description).To(Equal("newDescription"))
			Expect(newProject.RepositoryID.Valid).To(BeFalse())
			Expect(newProject.ShowIssues).To(BeTrue())
			Expect(newProject.ShowPullRequests).To(BeFalse())
		})
	})

	Describe("Repository", func() {
		Context("when repository exist", func() {
			It("should relate project to repository", func() {
				repositoryID := int64(12345)
				newProject, err := Create(uid, "new project", "description", repositoryID, "owner", "name", sql.NullString{})

				Expect(err).To(BeNil())
				repo, err := newProject.Repository()
				Expect(err).To(BeNil())
				Expect(repo.RepositoryID).To(Equal(repositoryID))
			})
		})
	})

	Describe("Lists", func() {
		var (
			newList, noneList *list.ListStruct
			newProject        *ProjectStruct
		)

		BeforeEach(func() {
			newProject, _ = Create(uid, "new project", "description", 0, "", "", sql.NullString{})
			newList = list.NewList(0, newProject.ID, newProject.UserID, "list title", "", sql.NullInt64{})
			_ = newList.Save(nil, nil)
			noneList = list.NewList(0, newProject.ID, newProject.UserID, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "", sql.NullInt64{})
			_ = noneList.Save(nil, nil)
		})
		It("should relate project and list", func() {
			lists := newProject.Lists()
			Expect(lists).NotTo(BeEmpty())
			Expect(lists[3].ID).To(Equal(newList.ID))
		})
		It("should not take none list", func() {
			lists := newProject.Lists()
			Expect(len(lists)).To(Equal(4))
		})

	})

	Describe("NoneList", func() {
		It("should contain only none list", func() {
			newProject, _ := Create(uid, "new project", "description", 0, "", "", sql.NullString{})
			noneList, err := newProject.NoneList()
			Expect(err).To(BeNil())
			Expect(noneList.Title.String).To(Equal(config.Element("init_list").(map[interface{}]interface{})["none"].(string)))
		})
	})
})
