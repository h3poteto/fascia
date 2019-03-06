package account

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/server/domains/project"
	domain "github.com/h3poteto/fascia/server/domains/user"
	repo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// InjectUserRepository inject db connection and return repository instance.
func InjectUserRepository() domain.Repository {
	return repo.New(InjectDB())
}

// RegistrationUser registers a user.
func RegistrationUser(email, password, passwordConfirm string) (*domain.User, error) {
	return domain.Registration(email, password, passwordConfirm, InjectUserRepository())
}

// FindUser finds a user.
func FindUser(id int64) (*domain.User, error) {
	return domain.Find(id, InjectUserRepository())
}

// FindUserByEmail finds a user.
func FindUserByEmail(email string) (*domain.User, error) {
	return domain.FindByEmail(email, InjectUserRepository())
}

// LoginUser check password and login.
func LoginUser(email, password string) (*domain.User, error) {
	return domain.Login(email, password, InjectUserRepository())
}

// FindOrCreateUserFromGithub creates a user from github.
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

func UserProjects(u *domain.User) ([]*project.Project, error) {
	return u.Projects(board.InjectProjectRepository())
}

/* Tests:
var _ = Describe("User", func() {
	Describe("FindOrCreateUserFromGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		Context("after registration from github", func() {
			user, err := FindOrCreateUserFromGithub(token)
			It("registration succeeded", func() {
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
