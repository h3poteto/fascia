package account

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	domain "github.com/h3poteto/fascia/server/domains/user"
	repo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func InjectUserRepository() domain.Repository {
	return repo.New(InjectDB())
}

func RegistrationUser(email, password, passwordConfirm string) (*domain.User, error) {
	return domain.Registration(email, password, passwordConfirm, InjectUserRepository())
}

func FindUser(id int64) (*domain.User, error) {
	return domain.Find(id, InjectUserRepository())
}

func FindUserByEmail(email string) (*domain.User, error) {
	return domain.FindByEmail(email, InjectUserRepository())
}

func LoginUser(email, password string) (*domain.User, error) {
	return domain.Login(email, password, InjectUserRepository())
}

func FindOrCreateUserFromGithub(token string) (*domain.User, error) {
	// GitHub authentication
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.Wrap(err, "github api error")
	}

	// TODO: Save not primary emails to login block.
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	emails, _, _ := client.Users.ListEmails(ctx, nil)
	var primaryEmail string
	for _, email := range emails {
		if *email.Primary {
			primaryEmail = *email.Email
		}
	}

	return domain.FindOrCreateFromGithub(githubUser, token, primaryEmail, InjectUserRepository())
}

/* Tests:
var _ = Describe("User", func() {
	Describe("FindOrCreateUserFromGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		Context("after registration from github", func() {
			user, err := FindOrCreateUserFromGithub(token)
			It("registration suceeded", func() {
				Expect(err).To(BeNil())
				Expect(user).NotTo(BeNil())
				findUser, _ := FindOrCreateUserFromGithub(token)
				Expect(findUser.UserEntity.ID).To(Equal(user.UserEntity.ID))
				Expect(findUser.UserEntity.ID).NotTo(BeZero())
			})
		})
		Context("after regist with email address", func() {
			email := "already_regist@example.com"
			RegistrationUser(email, "hogehoge", "hogehoge")
			user, _ := FindOrCreateUserFromGithub(token)
			It("should update github information", func() {
				Expect(user.UserEntity.OauthToken.Valid).To(BeTrue())
				Expect(user.UserEntity.OauthToken.String).To(Equal(token))
				Expect(user.UserEntity.UUID.Valid).To(BeTrue())
			})
		})
	})
})
*/
