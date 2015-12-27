package repository_test

import (
	"../db"
	"../project"
	. "../repository"
	"../user"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		newProject *project.ProjectStruct
		table      *sql.DB
	)
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table projects;")
		sql.Exec("truncate table repositories;")
		sql.Close()
	})
	JustBeforeEach(func() {
		email := "repository@example.com"
		password := "hogehoge"
		uid, _ := user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		newProject = project.NewProject(0, uid, "title", "desc")
		newProject.Save()
	})

	Describe("Save", func() {
		repositoryId := int64(123456)
		It("should create repository", func() {
			newRepository := NewRepository(0, newProject.Id, repositoryId, "owner", "repository_name")
			result := newRepository.Save()
			Expect(result).To(BeTrue())
		})
		It("should relate project to repository", func() {
			newRepository := NewRepository(0, newProject.Id, repositoryId, "owner", "repository_name")
			newRepository.Save()
			rows, _ := table.Query("select repositories.id from repositories inner join projects on repositories.project_id = projects.id;")

			var id int64
			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(id).To(Equal(newRepository.Id))
		})
	})
})
