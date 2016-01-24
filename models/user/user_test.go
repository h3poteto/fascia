package user_test

import (
	"../db"
	"../project"
	. "../user"
	"database/sql"
	"os"

	"github.com/google/go-github/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("User", func() {
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Exec("truncate table projects;")
		table.Close()
	})

	Describe("Registration", func() {
		email := "registration@example.com"
		password := "hogehoge"
		It("can regist", func() {
			id, err := Registration(email, password)
			Expect(err).To(BeNil())
			Expect(id).NotTo(Equal(int64(0)))
		})
		Context("after registration", func() {
			BeforeEach(func() {
				_, _ = Registration(email, password)
			})
			It("should save user in database", func() {
				mydb := &db.Database{}
				var database db.DB = mydb
				table := database.Init()
				rows, _ := table.Query("select id, email from users where email = ?;", email)

				var id int64
				var dbemail string
				for rows.Next() {
					err := rows.Scan(&id, &dbemail)
					if err != nil {
						panic(err.Error())
					}
				}
				Expect(dbemail).NotTo(Equal(""))
				Expect(id).NotTo(Equal(int64(0)))
			})
			It("cannot double regist", func() {
				id, err := Registration(email, password)
				Expect(err).NotTo(BeNil())
				Expect(id).To(Equal(int64(0)))
			})
		})

	})

	Describe("Login", func() {
		email := "login@example.com"
		password := "hogehoge"
		BeforeEach(func() {
			_, _ = Registration(email, password)
		})

		Context("when send correctly login information", func() {
			It("can login", func() {
				current_user, err := Login(email, password)
				Expect(err).To(BeNil())
				Expect(current_user.Email).To(Equal(email))
			})
		})
		Context("when send wrong login information", func() {
			It("cannot login", func() {
				current_user, err := Login(email, "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(current_user).To(BeNil())
			})
		})
		Context("when send wrong email address", func() {
			It("cannot login", func() {
				current_user, err := Login("hogehoge@example.com", password)
				Expect(err).NotTo(BeNil())
				Expect(current_user).To(BeNil())
			})
		})
		Context("when send wrong email address and password", func() {
			It("cannot login", func() {
				current_user, err := Login("hogehoge@example.com", "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(current_user).To(BeNil())
			})
		})
	})

	Describe("FindOrCreateGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		It("can regist through github", func() {
			_, err := FindOrCreateGithub(token)
			Expect(err).To(BeNil())
		})
		It("after regist through github, can search this user", func() {
			current_user, _ := FindOrCreateGithub(token)
			find_user, _ := FindOrCreateGithub(token)
			Expect(find_user.Id).To(Equal(current_user.Id))
			Expect(find_user.Id).NotTo(BeZero())
		})
		Context("after regist with email address", func() {
			email := "already_regist@example.com"
			var current_user *UserStruct
			BeforeEach(func() {
				Registration(email, "hogehoge")
				current_user, _ = FindOrCreateGithub(token)
			})
			It("should update github information", func() {
				Expect(current_user.OauthToken.Valid).To(BeTrue())
				Expect(current_user.OauthToken.String).To(Equal(token))
				Expect(current_user.Uuid.Valid).To(BeTrue())
			})
		})

	})

	Describe("Projects", func() {
		var (
			newProject   *project.ProjectStruct
			current_user *UserStruct
		)

		BeforeEach(func() {
			email := "project@example.com"
			password := "hogehoge"
			_, _ = Registration(email, password)
			mydb := &db.Database{}
			var database db.DB = mydb
			table := database.Init()
			rows, _ := table.Query("select id, email from users where email = ?;", email)

			var userid int64
			var dbemail string
			for rows.Next() {
				err := rows.Scan(&userid, &dbemail)
				if err != nil {
					panic(err.Error())
				}
			}
			newProject = project.NewProject(0, userid, "project title", "project desc", sql.NullInt64{})
			_ = newProject.Save()
			current_user = NewUser(userid, dbemail, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
		})
		It("ユーザとプロジェクトが関連づいていること", func() {
			projects := current_user.Projects()
			Expect(projects).NotTo(BeEmpty())
			Expect(projects[0].Id).To(Equal(newProject.Id))
		})
	})

	Describe("CreateGithubUser", func() {
		var result bool
		token := os.Getenv("TEST_TOKEN")
		BeforeEach(func() {
			newUser := NewUser(0, "", sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(oauth2.NoContext, ts)
			client := github.NewClient(tc)
			githubUser, _, _ := client.Users.Get("")
			result = newUser.CreateGithubUser(token, githubUser, "create_github_user@example.com")
		})
		It("ユーザが登録されること", func() {
			mydb := &db.Database{}
			var database db.DB = mydb
			table := database.Init()
			rows, err := table.Query("select id, oauth_token from users where oauth_token = ?;", token)
			if err != nil {
				panic(err.Error())
			}
			var id int64
			var oauthToken sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &oauthToken)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(result).To(BeTrue())
			Expect(oauthToken.Valid).To(BeTrue())
			Expect(id).NotTo(Equal(int64(0)))
		})
	})

	Describe("UpdateGithubUserInfo", func() {
		email := "update_github_user_info@example.com"
		token := os.Getenv("TEST_TOKEN")
		var current_user *UserStruct
		var result bool
		BeforeEach(func() {
			id, _ := Registration(email, "hogehoge")
			current_user, _ = CurrentUser(id)
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(oauth2.NoContext, ts)
			client := github.NewClient(tc)
			githubUser, _, _ := client.Users.Get("")
			result = current_user.UpdateGithubUserInfo(token, githubUser)
		})
		It("ユーザ情報がアップデートされること", func() {
			mydb := &db.Database{}
			var database db.DB = mydb
			table := database.Init()
			rows, err := table.Query("select id, uuid, oauth_token from users where email = ?;", email)
			if err != nil {
				panic(err.Error())
			}
			var id int64
			var oauthToken sql.NullString
			var uuid sql.NullInt64
			for rows.Next() {
				err := rows.Scan(&id, &uuid, &oauthToken)
				if err != nil {
					panic(err.Error())
				}
			}
			Expect(result).To(BeTrue())
			Expect(oauthToken.Valid).To(BeTrue())
			Expect(oauthToken.String).To(Equal(token))
			Expect(uuid.Valid).To(BeTrue())
			Expect(id).NotTo(Equal(int64(0)))
		})
	})
})
