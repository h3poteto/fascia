package user_test

import (
	seed "../../db/seed"
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
	BeforeEach(func() {
		seed.ListOptions()
	})
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Exec("truncate table projects;")
		table.Exec("truncate table list_options;")
		table.Close()
	})

	Describe("Validation", func() {
		Context("when email is empty", func() {
			email := ""
			password := "hogehoge"
			It("should validation failed", func() {
				result := Validation(email, password, password)
				Expect(result).To(BeFalse())
			})
		})
		Context("when email does not contain domain", func() {
			email := "abc@def"
			password := "hogehoge"
			It("should validation failed", func() {
				result := Validation(email, password, password)
				Expect(result).To(BeFalse())
			})
		})
		Context("when passwod length is less than 8 characters", func() {
			email := "abc@def.com"
			password := "hoho"
			It("should validation failed", func() {
				result := Validation(email, password, password)
				Expect(result).To(BeFalse())
			})
		})
		Context("when password is empty", func() {
			email := "abc@def.com"
			password := ""
			It("should validation failed", func() {
				result := Validation(email, password, password)
				Expect(result).To(BeFalse())
			})
		})
		Context("when password disagree with password comfirm", func() {
			email := "abc@def.com"
			password := "hogehoge"
			passwordConfirm := "fugafuga"
			It("should validation failed", func() {
				result := Validation(email, password, passwordConfirm)
				Expect(result).To(BeFalse())
			})
		})
		Context("when email and password is right", func() {
			email := "abc@def.com"
			password := "hogehoge"
			It("should validation success", func() {
				result := Validation(email, password, password)
				Expect(result).To(BeTrue())
			})
		})
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
						panic(err)
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
				currentUser, err := Login(email, password)
				Expect(err).To(BeNil())
				Expect(currentUser.Email).To(Equal(email))
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

	Describe("FindOrCreateGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		It("can regist through github", func() {
			_, err := FindOrCreateGithub(token)
			Expect(err).To(BeNil())
		})
		It("after regist through github, can search this user", func() {
			currentUser, _ := FindOrCreateGithub(token)
			find_user, _ := FindOrCreateGithub(token)
			Expect(find_user.ID).To(Equal(currentUser.ID))
			Expect(find_user.ID).NotTo(BeZero())
		})
		Context("after regist with email address", func() {
			email := "already_regist@example.com"
			var currentUser *UserStruct
			BeforeEach(func() {
				Registration(email, "hogehoge")
				currentUser, _ = FindOrCreateGithub(token)
			})
			It("should update github information", func() {
				Expect(currentUser.OauthToken.Valid).To(BeTrue())
				Expect(currentUser.OauthToken.String).To(Equal(token))
				Expect(currentUser.Uuid.Valid).To(BeTrue())
			})
		})

	})

	Describe("Projects", func() {
		var (
			newProject  *project.ProjectStruct
			currentUser *UserStruct
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
					panic(err)
				}
			}
			var err error
			newProject, err = project.Create(userid, "title", "desc", 0, "", "", sql.NullString{})
			if err != nil {
				panic(err)
			}
			currentUser = NewUser(userid, dbemail, sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{})
		})
		It("ユーザとプロジェクトが関連づいていること", func() {
			projects := currentUser.Projects()
			Expect(projects).NotTo(BeEmpty())
			Expect(projects[0].ID).To(Equal(newProject.ID))
		})
	})

	Describe("CreateGithubUser", func() {
		var result error
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
			Expect(err).To(BeNil())
			var id int64
			var oauthToken sql.NullString
			for rows.Next() {
				err := rows.Scan(&id, &oauthToken)
				if err != nil {
					panic(err)
				}
			}
			Expect(result).To(BeNil())
			Expect(oauthToken.Valid).To(BeTrue())
			Expect(id).NotTo(Equal(int64(0)))
		})
	})

	Describe("UpdateGithubUserInfo", func() {
		email := "update_github_user_info@example.com"
		token := os.Getenv("TEST_TOKEN")
		var currentUser *UserStruct
		var result error
		BeforeEach(func() {
			id, _ := Registration(email, "hogehoge")
			currentUser, _ = CurrentUser(id)
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(oauth2.NoContext, ts)
			client := github.NewClient(tc)
			githubUser, _, _ := client.Users.Get("")
			result = currentUser.UpdateGithubUserInfo(token, githubUser)
		})
		It("ユーザ情報がアップデートされること", func() {
			mydb := &db.Database{}
			var database db.DB = mydb
			table := database.Init()
			rows, err := table.Query("select id, uuid, oauth_token from users where email = ?;", email)
			Expect(err).To(BeNil())
			var id int64
			var oauthToken sql.NullString
			var uuid sql.NullInt64
			for rows.Next() {
				err := rows.Scan(&id, &uuid, &oauthToken)
				if err != nil {
					panic(err)
				}
			}
			Expect(result).To(BeNil())
			Expect(oauthToken.Valid).To(BeTrue())
			Expect(oauthToken.String).To(Equal(token))
			Expect(uuid.Valid).To(BeTrue())
			Expect(id).NotTo(Equal(int64(0)))
		})
	})
})
