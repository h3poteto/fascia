package user_test

import (
	"os"
	"database/sql"
	"../project"
	. "../user"
	"../db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

var _ = Describe("User", func() {
	var currentdb string

	BeforeEach(func() {
		testdb := os.Getenv("DB_TEST_NAME")
		currentdb = os.Getenv("DB_NAME")
		os.Setenv("DB_NAME", testdb)
	})
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Exec("truncate table projects;")
		table.Close()
		os.Setenv("DB_NAME", currentdb)
	})

	Describe("Registration", func() {
		email := "registration@example.com"
		password := "hogehoge"
		It("登録できること", func() {
			id, reg := Registration(email, password)
			Expect(reg).To(BeTrue())
			Expect(id).NotTo(Equal(int64(0)))
		})
		Context("登録後", func() {
			BeforeEach(func() {
				_, _ = Registration(email, password)
			})
			It("DBにユーザが保存されていること", func() {
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
			It("ユーザが二重登録できないこと", func() {
				id, reg := Registration(email, password)
				Expect(reg).To(BeFalse())
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

		Context("正しいログイン情報のとき", func() {
			It("ログインできること", func() {
				current_user, err := Login(email, password)
				Expect(err).To(BeNil())
				Expect(current_user.Email).To(Equal(email))
			})
		})
		Context("パスワードを間違えているとき", func() {
			It("ログインできないこと", func() {
				current_user, err := Login(email, "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(current_user.Email).NotTo(Equal(email))
			})
		})
		Context("メールアドレスを間違えているとき", func() {
			It("ログインできないこと", func() {
				current_user, err := Login("hogehoge@example.com", password)
				Expect(err).NotTo(BeNil())
				Expect(current_user.Email).NotTo(Equal(email))
			})
		})
		Context("メールアドレスもパスワードも間違えているとき", func() {
			It("ログインできないこと", func() {
				current_user, err := Login("hogehoge@example.com", "fugafuga")
				Expect(err).NotTo(BeNil())
				Expect(current_user.Email).NotTo(Equal(email))
			})
		})
	})

	Describe("FindOrCreateGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		It("Github経由で新規登録できること", func() {
			_ , err := FindOrCreateGithub(token)
			Expect(err).To(BeNil())
		})
		It("github登録後であればすでに登録されているユーザを探せること", func() {
			current_user, _ := FindOrCreateGithub(token)
			find_user, _ := FindOrCreateGithub(token)
			Expect(find_user.Id).To(Equal(current_user.Id))
			Expect(find_user.Id).NotTo(BeZero())
		})
		Context("ユーザがEmailで登録済みだったとき", func() {
			email := "already_regist@example.com"
			var current_user *UserStruct
			BeforeEach(func() {
				Registration(email, "hogehoge")
				current_user, _ = FindOrCreateGithub(token)
			})
			It("github情報が更新されること", func() {
				Expect(current_user.OauthToken.Valid).To(BeTrue())
				Expect(current_user.OauthToken.String).To(Equal(token))
				Expect(current_user.Uuid.Valid).To(BeTrue())
			})
		})

	})

	Describe("Projects", func() {
		var (
			newProject *project.ProjectStruct
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
			newProject = project.NewProject(0, userid, "project title")
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
