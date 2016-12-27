package project_test

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/entities/project"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	var (
		newProject *Project
		uid        int64
		database   *sql.DB
	)

	BeforeEach(func() {
		seed.Seeds()
		email := "save@example.com"
		password := "hogehoge"
		user, err := handlers.RegistrationUser(email, password, password)
		if err != nil {
			panic(err)
		}
		uid = user.UserEntity.UserModel.ID
		database = db.SharedInstance().Connection
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
		database.Exec("truncate table repositories;")
		database.Exec("truncate table lists;")
		database.Exec("truncate table list_options;")
	})

	Describe("Update", func() {
		BeforeEach(func() {
			var repoID sql.NullInt64
			newProject = New(0, uid, "new project", "description", repoID, false, false)
			newProject.Save(nil)
		})
		It("should set new value", func() {
			err := newProject.Update("newTitle", "newDescription", true, false)
			Expect(err).To(BeNil())
			Expect(newProject.ProjectModel.Title).To(Equal("newTitle"))
			Expect(newProject.ProjectModel.Description).To(Equal("newDescription"))
			Expect(newProject.ProjectModel.RepositoryID.Valid).To(BeFalse())
			Expect(newProject.ProjectModel.ShowIssues).To(BeTrue())
			Expect(newProject.ProjectModel.ShowPullRequests).To(BeFalse())
		})
	})

	Describe("CreateInitialLists", func() {
		var (
			tx         *sql.Tx
			newProject *Project
		)
		BeforeEach(func() {
			tx, _ = db.SharedInstance().Connection.Begin()
			newProject = New(0, uid, "new project", "description", sql.NullInt64{}, false, false)
			newProject.Save(tx)
		})
		It("should success to create", func() {
			err := newProject.CreateInitialLists(tx)
			Expect(err).To(BeNil())
			err = tx.Commit()
			Expect(err).To(BeNil())
		})
	})

	Describe("Lists and NoneLists", func() {
		var (
			newProject *Project
		)
		BeforeEach(func() {
			tx, _ := db.SharedInstance().Connection.Begin()
			newProject = New(0, uid, "new project", "description", sql.NullInt64{}, false, false)
			newProject.Save(tx)
			newProject.CreateInitialLists(tx)
			tx.Commit()
		})
		Describe("Lists", func() {
			It("should relate project and list", func() {
				lists, err := newProject.Lists()
				Expect(err).To(BeNil())
				Expect(lists).NotTo(BeEmpty())
			})
		})

		Describe("NoneList", func() {
			It("should contain only none list", func() {
				noneList, err := newProject.NoneList()
				Expect(err).To(BeNil())
				Expect(noneList.ListModel.Title.String).To(Equal(config.Element("init_list").(map[interface{}]interface{})["none"].(string)))
			})
		})
	})
})
