package controllers_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/google/go-github/github"
	. "github.com/h3poteto/fascia/server"
	"github.com/h3poteto/fascia/server/models/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

var _ = Describe("GithubController", func() {
	var (
		ts       *httptest.Server
		database *sql.DB
	)
	userEmail := "github@example.com"
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
	})
	AfterEach(func() {
		ts.Close()
		database.Exec("truncate table users;")
	})
	JustBeforeEach(func() {
		LoginFaker(ts, userEmail, "hogehoge")
		// Oauthのログインテストはリダイレクトまでしか実行できないため，OauthTokenは偽装しておくしかない
		database = db.SharedInstance().Connection

		token := os.Getenv("TEST_TOKEN")
		// github認証
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		githubUser, _, err := client.Users.Get("")
		if err != nil {
			log.Fatal(err)
		}
		database.Exec("update users set provider = ?, oauth_token = ?, user_name = ?, uuid = ?, avatar_url = ? where email = ?;", "github", token, *githubUser.Login, *githubUser.ID, *githubUser.AvatarURL, userEmail)

	})
	Describe("Repositories", func() {
		It("should receive repositories", func() {
			res, err := http.Get(ts.URL + "/github/repositories")
			Expect(err).To(BeNil())
			contents, status := ParseJson(res)
			Expect(status).To(Equal(http.StatusOK))
			parseContents := contents.([]interface{})
			Expect(parseContents[0]).NotTo(BeNil())
		})
	})

})
