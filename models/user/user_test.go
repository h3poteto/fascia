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
			reg := Registration(email, password)
			Expect(reg).To(BeTrue())
		})
		Context("登録後", func() {
			BeforeEach(func() {
				_ = Registration(email, password)
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
				reg := Registration(email, password)
				Expect(reg).To(BeFalse())
			})
		})

	})

	Describe("Login", func() {
		email := "login@example.com"
		password := "hogehoge"
		BeforeEach(func() {
			_ = Registration(email, password)
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
		It("登録後であればすでに登録されているユーザを探せること", func() {
			current_user, _ := FindOrCreateGithub(token)
			find_user, _ := FindOrCreateGithub(token)
			Expect(find_user.Id).To(Equal(current_user.Id))
			Expect(find_user.Id).NotTo(BeZero())
		})

	})

	Describe("Project", func() {
		var (
			newProject *project.ProjectStruct
			current_user *UserStruct
		)

		BeforeEach(func() {
			email := "project@example.com"
			password := "hogehoge"
			_ = Registration(email, password)
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
})
