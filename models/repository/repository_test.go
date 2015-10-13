package repository_test

import (
	"os"
	"database/sql"
	"../db"
	. "../repository"
	"../project"
	"../user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		newProject *project.ProjectStruct
		currentdb string
		table *sql.DB
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
		sql.Exec("truncate table repositories;")
		sql.Close()
		os.Setenv("DB_NAME", currentdb)
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
		It("リポジトリが新規作成できること", func() {
			newRepository := NewRepository(0, newProject.Id, repositoryId, "title")
			result := newRepository.Save()
			Expect(result).To(BeTrue())
		})
		It("リポジトリとプロジェクトが関連づくこと", func() {
			newRepository := NewRepository(0, newProject.Id, repositoryId, "title")
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
