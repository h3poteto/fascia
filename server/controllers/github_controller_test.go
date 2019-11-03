package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/database"
	. "github.com/h3poteto/fascia/server/controllers"
	userRepo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("GithubController", func() {
	var (
		e   *echo.Echo
		rec *httptest.ResponseRecorder
		db  *sql.DB
	)
	userEmail := "github@example.com"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	// Oauthのログインテストはリダイレクトまでしか実行できないため，OauthTokenは偽装しておくしかない
	BeforeEach(func() {
		db = database.SharedInstance().Connection

		token := os.Getenv("TEST_TOKEN")
		// github認証
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		client := github.NewClient(tc)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		githubUser, _, err := client.Users.Get(ctx, "")
		if err != nil {
			log.Fatal(err)
		}
		repo := userRepo.New(db)
		repo.Create(
			userEmail,
			"hogehoge",
			sql.NullString{String: "github", Valid: true},
			sql.NullString{String: token, Valid: true},
			sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true},
			sql.NullString{String: *githubUser.Login, Valid: true},
			sql.NullString{String: *githubUser.AvatarURL, Valid: true})

	})
	Describe("Repositories", func() {
		It("should receive repositories", func() {
			u, _ := account.FindUserByEmail(userEmail)
			req := httptest.NewRequest(echo.GET, "/github/repositories", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, u)
			resource := Github{}
			err := resource.Repositories(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			parseContents := contents.([]interface{})
			Expect(parseContents[0]).NotTo(BeNil())
		})
	})

})
