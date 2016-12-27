package services_test

import (
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"
	. "github.com/h3poteto/fascia/server/services"

	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectService", func() {
	var (
		uid      int64
		database *sql.DB
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
		database.Exec("truncate table lists;")
		database.Exec("truncate table list_options;")
	})

	Describe("Create", func() {
		projectService := NewProject(nil)
		Context("when did not set repositoryID", func() {
			It("should create new project", func() {
				newProject, err := projectService.Create(uid, "new project", "description", 0, sql.NullString{})
				Expect(err).To(BeNil())
				lists, _ := newProject.Lists()
				Expect(len(lists)).To(Equal(3))
				Expect(newProject.NoneList()).NotTo(BeNil())
				Expect(newProject.ProjectModel.ShowIssues).To(BeTrue())
				Expect(newProject.ProjectModel.ShowPullRequests).To(BeTrue())
			})
			It("should relate user and project", func() {
				newProject, _ := projectService.Create(uid, "new project", "description", 0, sql.NullString{})
				rows, _ := database.Query("select id, user_id, title, description from projects where id = ?;", newProject.ProjectModel.ID)

				var id int64
				var userID sql.NullInt64
				var title string
				var description string

				for rows.Next() {
					err := rows.Scan(&id, &userID, &title, &description)
					if err != nil {
						panic(err)
					}
				}
				Expect(userID.Valid).To(BeTrue())
				Expect(userID.Int64).To(Equal(uid))
			})
		})
	})
})
