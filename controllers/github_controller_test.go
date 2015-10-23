package controllers_test

import (
	"os"
	"net/http"
	"net/http/httptest"
	. "../../fascia"
	"../models/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

var _ = Describe("GithubController", func() {
	var (
		ts *httptest.Server
		currentdb string
	)
	userEmail := "github@example.com"
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
		testdb := os.Getenv("DB_TEST_NAME")
		currentdb = os.Getenv("DB_NAME")
		os.Setenv("DB_NAME", testdb)
	})
	AfterEach(func() {
		ts.Close()
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Close()
		os.Setenv("DB_NAME", currentdb)
	})
	JustBeforeEach(func() {
		LoginFaker(ts, userEmail, "hogehoge")
		// Oauthのログインテストはリダイレクトまでしか実行できないため，OauthTokenは偽装しておくしかない
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()

		token := os.Getenv("TEST_TOKEN")
		// github認証
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		githubUser, _, _ := client.Users.Get("")

		table.Exec("update users set provider = ?, oauth_token =?, user_name = ?, uuid = ?, avatar_url = ? where email = ?;", "github", token, *githubUser.Login, *githubUser.ID, *githubUser.AvatarURL, userEmail)

	})
	Describe("Repositories", func() {
		It("リポジトリが取得できること", func() {
			res, err := http.Get(ts.URL + "/github/repositories")
			Expect(err).To(BeNil())
			contents, status := ParseJson(res)
			Expect(status).To(Equal(http.StatusOK))
			parseContents := contents.([]interface{})
			Expect(parseContents[0]).NotTo(BeNil())
		})
	})

})
