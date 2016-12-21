package repository

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/server/models/repository"

	"github.com/pkg/errors"
)

// Repository has repository model object
type Repository struct {
	RepositoryModel *repository.Repository
}

// New returns a repository entity
func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *Repository {
	return &Repository{
		RepositoryModel: repository.New(id, repositoryID, owner, name, webhookKey),
	}
}

// FindByGithubRepoID find repository entity according to repository id in github
func FindByGithubRepoID(id int64) (*Repository, error) {
	r, err := repository.FindByGithubRepoID(id)
	if err != nil {
		return nil, err
	}
	return &Repository{
		RepositoryModel: r,
	}, nil
}

// CreateRepository create repository record based on github repository
func CreateRepository(ID int, oauthToken string) (*Repository, error) {
	// confirm github
	h := hub.New(oauthToken)
	githubRepo, err := h.GetRepository(ID)
	if err != nil {
		return nil, err
	}
	// generate webhook key
	key := GenerateWebhookKey(*githubRepo.Name)
	// save
	repo := New(0, int64(*githubRepo.ID), *githubRepo.Owner.Login, *githubRepo.Name, key)
	err = repo.Save()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// GenerateWebhookKey create new md5 hash
func GenerateWebhookKey(seed string) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

// Save call repository model save
func (r *Repository) Save() error {
	return r.RepositoryModel.Save()
}

// Authenticate is check token and webhookKey with response
func (r *Repository) Authenticate(token string, response []byte) error {
	mac := hmac.New(sha1.New, []byte(r.RepositoryModel.WebhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}
