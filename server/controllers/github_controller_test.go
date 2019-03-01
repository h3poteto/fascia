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
	"github.com/h3poteto/fascia/server/domains/user"
	usecase "github.com/h3poteto/fascia/server/usecases/account"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("GithubController", func() {
	var (
		e    *echo.Echo
		rec  *httptest.ResponseRecorder
		db   *sql.DB
		user *user.User
	)
	userEmail := "github@example.com"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	// Oauthのログインテストはリダイレクトまでしか実行できないため，OauthTokenは偽装しておくしかない
	JustBeforeEach(func() {
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
		user, _ = usecase.RegistrationUser(userEmail, "hogehoge", "hogehoge")
		db.Exec("update users set provider = ?, oauth_token = ?, user_name = ?, uuid = ?, avatar_url = ? where email = ?;", "github", token, *githubUser.Login, *githubUser.ID, *githubUser.AvatarURL, userEmail)

	})
	Describe("Repositories", func() {
		It("should receive repositories", func() {
			user, _ = usecase.FindUserByEmail(userEmail)
			req := httptest.NewRequest(echo.GET, "/github/repositories", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, userEmail, "hogehoge")
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
