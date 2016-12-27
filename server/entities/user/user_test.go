package user_test

import (
	"database/sql"

	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/entities/user"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	var (
		database *sql.DB
	)
	BeforeEach(func() {
		seed.Seeds()
		database = db.SharedInstance().Connection
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
		database.Exec("truncate table list_options;")
	})

	Describe("Registration", func() {
		email := "registration@example.com"
		password := "hogehoge"
		It("can regist", func() {
			user, err := Registration(email, password, password)
			Expect(err).To(BeNil())
			Expect(user.UserModel.ID).NotTo(Equal(int64(0)))
		})
		Context("after registration", func() {
			BeforeEach(func() {
				Registration(email, password, password)
			})
			It("should save user in database", func() {
				rows, _ := database.Query("select id, email from users where email = ?;", email)

				var id int64
				var dbemail string
				for rows.Next() {
					err := rows.Scan(&id, &dbemail)
					if err != nil {
						panic(err)
					}
				}
				Expect(dbemail).NotTo(Equal(""))
				Expect(id).NotTo(Equal(int64(0)))
			})
			It("cannot double regist", func() {
				user, err := Registration(email, password, password)
				Expect(err).NotTo(BeNil())
				Expect(user).To(BeNil())
			})
		})

	})

	Describe("Login", func() {
		email := "login@example.com"
		password := "hogehoge"
		BeforeEach(func() {
			Registration(email, password, password)
		})

		Context("when send correctly login information", func() {
			It("can login", func() {
				currentUser, err := Login(email, password)
				Expect(err).To(BeNil())
				Expect(currentUser.UserModel.Email).To(Equal(email))
			})
		})
		Context("when send wrong login information", func() {
			It("cannot login", func() {
				currentUser, err := Login(email, "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(currentUser).To(BeNil())
			})
		})
		Context("when send wrong email address", func() {
			It("cannot login", func() {
				currentUser, err := Login("hogehoge@example.com", password)
				Expect(err).NotTo(BeNil())
				Expect(currentUser).To(BeNil())
			})
		})
		Context("when send wrong email address and password", func() {
			It("cannot login", func() {
				currentUser, err := Login("hogehoge@example.com", "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(currentUser).To(BeNil())
			})
		})
	})

	Describe("Projects", func() {
		var (
			newProject  *services.Project
			currentUser *User
		)

		BeforeEach(func() {
			email := "project@example.com"
			password := "hogehoge"
			Registration(email, password, password)
			rows, _ := database.Query("select id, email from users where email = ?;", email)

			var userid int64
			var dbemail string
			for rows.Next() {
				err := rows.Scan(&userid, &dbemail)
				if err != nil {
					panic(err)
				}
			}
			var err error
			newProject, err = handlers.CreateProject(userid, "title", "desc", 0, sql.NullString{})
			if err != nil {
				panic(err)
			}
			currentUser = New(userid, dbemail, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
		})
		It("ユーザとプロジェクトが関連づいていること", func() {
			projects, err := currentUser.Projects()
			Expect(err).To(BeNil())
			Expect(projects).NotTo(BeEmpty())
			Expect(projects[0].ProjectModel.ID).To(Equal(newProject.ProjectEntity.ProjectModel.ID))
		})
	})
})
