package account

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/domains/project"
	domain "github.com/h3poteto/fascia/server/domains/user"
	repo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// InjectDB set DB connection from connection pool.
func InjectDB() *sql.DB {
	return database.SharedInstance().Connection
}

// InjectUserRepository inject db connection and return repository instance.
func InjectUserRepository() domain.Repository {
	return repo.New(InjectDB())
}

// FindUser finds a user.
func FindUser(id int64) (*domain.User, error) {
	return domain.Find(id, InjectUserRepository())
}

// FindUserByEmail finds a user.
func FindUserByEmail(email string) (*domain.User, error) {
	return domain.FindByEmail(email, InjectUserRepository())
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

// UserProjects returns projects related the user.
func UserProjects(u *domain.User) ([]*project.Project, error) {
	return u.Projects(board.InjectProjectRepository())
}
