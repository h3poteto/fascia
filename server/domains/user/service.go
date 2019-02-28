package user

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

// hashPassword generate hash password
func hashPassword(password string) ([]byte, error) {
	bytePassword := []byte(password)
	cost := 10
	hashed, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate password")
	}
	err = bcrypt.CompareHashAndPassword(hashed, bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "did not match password")
	}
	return hashed, nil
}

// Find returns a user entity
func Find(targetID int64, infrastructure Repository) (*User, error) {
	id, email, password, provider, oauthToken, uuid, userName, avatar, err := infrastructure.Find(targetID)
	if err != nil {
		return nil, err
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatar, infrastructure), nil

}

// FindByEmail returns a user entity
func FindByEmail(targetEmail string, infrastructure Repository) (*User, error) {
	id, email, password, provider, oauthToken, uuid, userName, avatar, err := infrastructure.FindByEmail(targetEmail)
	if err != nil {
		return nil, err
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatar, infrastructure), nil
}

// Registration create a new user record
func Registration(email, password, passwordConfirm string, infrastructure Repository) (*User, error) {
	if password != passwordConfirm {
		return nil, errors.New("password is incorrect")
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	u := New(0, email, string(hashed), sql.NullString{}, sql.NullString{}, sql.NullInt64{}, sql.NullString{}, sql.NullString{}, infrastructure)
	if err := u.create(); err != nil {
		return nil, err
	}
	return u, nil
}

// Login authenticate email and password
func Login(targetEmail string, targetPassword string, infrastructure Repository) (*User, error) {
	id, email, password, provider, oauthToken, uuid, userName, avatar, err := infrastructure.FindByEmail(targetEmail)
	if err != nil {
		return nil, err
	}
	bytePassword := []byte(targetPassword)
	err = bcrypt.CompareHashAndPassword([]byte(password), bytePassword)
	if err != nil {
		return nil, errors.Wrap(err, "password did not match")
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatar, infrastructure), nil
}

// FindOrCreateFromGithub create or update user based on github user.
func FindOrCreateFromGithub(githubUser *github.User, token string, primaryEmail string, infrastructure Repository) (*User, error) {
	id, email, password, provider, oauthToken, uuid, userName, avatar, err := infrastructure.FindByEmail(primaryEmail)
	if err != nil {
		// Create new user if does not exist
		u, err := createGithubUser(githubUser, token, primaryEmail, infrastructure)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	// When oauth information is updated, we have to update user
	if !oauthToken.Valid || oauthToken.String != token {
		if err := updateGithubUser(githubUser, id, email, token, infrastructure); err != nil {
			return nil, err
		}
	}
	return New(id, email, password, provider, oauthToken, uuid, userName, avatar, infrastructure), nil
}

func createGithubUser(githubUser *github.User, token, primaryEmail string, infrastructure Repository) (*User, error) {
	bytePassword, err := hashPassword(randomString())
	if err != nil {
		return nil, err
	}
	provider := sql.NullString{String: "github", Valid: true}
	oauthToken := sql.NullString{String: token, Valid: true}
	userName := sql.NullString{String: *githubUser.Login, Valid: true}
	uuid := sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	avatar := sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	u := New(0, primaryEmail, string(bytePassword), provider, oauthToken, uuid, userName, avatar, infrastructure)
	if err := u.create(); err != nil {
		return nil, err
	}
	return u, nil
}

func updateGithubUser(githubUser *github.User, id int64, email, token string, infrastructure Repository) error {
	provider := sql.NullString{String: "github", Valid: true}
	oauthToken := sql.NullString{String: token, Valid: true}
	userName := sql.NullString{String: *githubUser.Login, Valid: true}
	uuid := sql.NullInt64{Int64: int64(*githubUser.ID), Valid: true}
	avatar := sql.NullString{String: *githubUser.AvatarURL, Valid: true}
	u := New(id, email, "", provider, oauthToken, uuid, userName, avatar, infrastructure)
	return u.Update()
}
