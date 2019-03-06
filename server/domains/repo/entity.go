package repo

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"

	"github.com/pkg/errors"
)

// Repository has repository model object
type Repo struct {
	ID             int64
	RepositoryID   int64
	Owner          sql.NullString
	Name           sql.NullString
	WebhookKey     string
	infrastructure Repository
}

type Repository interface {
	FindByGithubRepoID(int64) (int64, int64, sql.NullString, sql.NullString, string, error)
	FindByProjectID(int64) (int64, int64, sql.NullString, sql.NullString, string, error)
	Create(int64, sql.NullString, sql.NullString, string) (int64, error)
}

// New returns a repository entity
func New(id int64, repositoryID int64, owner sql.NullString, name sql.NullString, webhookKey string, infrastructure Repository) *Repo {
	return &Repo{
		id,
		repositoryID,
		owner,
		name,
		webhookKey,
		infrastructure,
	}
}

func (r *Repo) Create() error {
	id, err := r.infrastructure.Create(r.RepositoryID, r.Owner, r.Name, r.WebhookKey)
	if err != nil {
		return err
	}
	r.ID = id
	return nil
}

// Authenticate is check token and webhookKey with response
func (r *Repo) Authenticate(token string, response []byte) error {
	mac := hmac.New(sha1.New, []byte(r.WebhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}
