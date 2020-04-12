package account

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"github.com/h3poteto/fascia/lib/modules/database"
	"github.com/h3poteto/fascia/server/domains/project"
	domain "github.com/h3poteto/fascia/server/domains/user"
	repo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	repo := InjectUserRepository()
	return repo.Find(id)
}

// FindUserByEmail finds a user.
func FindUserByEmail(email string) (*domain.User, error) {
	repo := InjectUserRepository()
	return repo.FindByEmail(email)
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

	return findOrCreateFromGithub(githubUser, token, primaryEmail)
}

// UserProjects returns projects related the user.
func UserProjects(u *domain.User) ([]*project.Project, error) {
	infra := board.InjectProjectRepository()
	return infra.Projects(u.ID)
}

// FindOrCreateFromGithub create or update user based on github user.
func findOrCreateFromGithub(githubUser *github.User, token string, primaryEmail string) (*domain.User, error) {
	repo := InjectUserRepository()
	u, err := repo.FindByEmail(primaryEmail)
	if err != nil {
		// Create new user if does not exist
		u, err := createGithubUser(githubUser, token, primaryEmail, repo)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	// When oauth information is updated, we have to update user
	if !u.OauthToken.Valid || u.OauthToken.String != token {
		u.UpdateGithubUser(githubUser, u.ID, u.Email, token)
		if err := repo.Update(u.ID, u.Email, u.Provider, u.OauthToken, u.UUID, u.UserName, u.Avatar); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func createGithubUser(githubUser *github.User, token, primaryEmail string, repo domain.Repository) (*domain.User, error) {
	bytePassword, err := hashPassword(randomString())
	if err != nil {
		return nil, err
	}
	provider := sql.NullString{String: "github", Valid: true}
	oauthToken := sql.NullString{String: token, Valid: true}
	userName := sql.NullString{String: *githubUser.Login, Valid: true}
	uuid := sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	avatar := sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	id, err := repo.Create(primaryEmail, string(bytePassword), provider, oauthToken, uuid, userName, avatar)
	if err != nil {
		return nil, err
	}
	u, err := repo.Find(id)
	return u, nil
}

// hashPassword generate hash password
func hashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashed, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		return nil, errors.Wrap(err, "generate password error")
	}
	err = bcrypt.CompareHashAndPassword(hashed, bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "compare password error")
	}
	return hashed, nil
}

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

// UpdatePassword updates the user password.
func UpdatePassword(id int64, password, passwordConfirm string) (*domain.User, error) {
	user, err := FindUser(id)
	if err != nil {
		return nil, err
	}
	if password != passwordConfirm {
		return nil, errors.New("password is not matched")
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Update(user.Email, string(hashed), user.Provider, user.OauthToken, user.UUID, user.UserName, user.Avatar)

	repo := InjectUserRepository()
	if err := repo.UpdatePassword(user.ID, user.HashedPassword); err != nil {
		return nil, err
	}
	return user, nil

}
