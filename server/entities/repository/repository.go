package repository

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/server/infrastructures/repository"

	"github.com/pkg/errors"
)

// Repository has repository model object
type Repository struct {
	ID             int64
	RepositoryID   int64
	Owner          sql.NullString
	Name           sql.NullString
	WebhookKey     string
	infrastructure *repository.Repository
}

// New returns a repository entity
func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *Repository {
	infrastructure := repository.New(id, repositoryID, owner, name, webhookKey)
	r := &Repository{
		infrastructure: infrastructure,
	}
	r.reload()
	return r
}

func (r *Repository) reflect() {
	r.infrastructure.ID = r.ID
	r.infrastructure.RepositoryID = r.RepositoryID
	r.infrastructure.Owner = r.Owner
	r.infrastructure.Name = r.Name
	r.infrastructure.WebhookKey = r.WebhookKey
}

func (r *Repository) reload() error {
	if r.RepositoryID != 0 {
		latestRepository, err := repository.FindByGithubRepoID(r.RepositoryID)
		if err != nil {
			return err
		}
		r.infrastructure = latestRepository
	}
	r.ID = r.infrastructure.ID
	r.RepositoryID = r.infrastructure.RepositoryID
	r.Owner = r.infrastructure.Owner
	r.Name = r.infrastructure.Name
	r.WebhookKey = r.infrastructure.WebhookKey
	return nil
}

// CreateRepository create repository record based on github repository
func CreateRepository(id int64, oauthToken string) (*Repository, error) {
	// confirm github
	h := hub.New(oauthToken)
	githubRepo, err := h.GetRepository(int(id))
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
	r.reflect()
	return r.infrastructure.Save()
}

// Authenticate is check token and webhookKey with response
func (r *Repository) Authenticate(token string, response []byte) error {
	mac := hmac.New(sha1.New, []byte(r.infrastructure.WebhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}
