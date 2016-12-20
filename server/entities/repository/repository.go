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

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type Repository struct {
	RepositoryModel *repository.RepositoryStruct
}

func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *Repository {
	return &Repository{
		RepositoryModel: repository.New(id, repositoryID, owner, name, webhookKey),
	}
}

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

func (r *Repository) CheckLabelPresent(token, title string) (*github.Label, error) {
	return hub.CheckLabelPresent(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title)
}

func (r *Repository) CreateGithubLabel(token, title, color string) (*github.Label, error) {
	return hub.CreateGithubLabel(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, color)
}

func (r *Repository) UpdateGithubLabel(token, originalTitle, title, color string) (*github.Label, error) {
	return hub.UpdateGithubLabel(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, originalTitle, title, color)
}

func (r *Repository) CreateGithubIssue(token, title, description string, labels []string) (*github.Issue, error) {
	return hub.CreateGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, description, labels)
}

func (r *Repository) EditGithubIssue(token, title, description, state string, issueNumber int, labels []string) (bool, error) {
	return hub.EditGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, title, description, state, issueNumber, labels)
}

func (r *Repository) GetGithubIssue(token string, number int) (*github.Issue, error) {
	return hub.GetGithubIssue(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, number)
}

func (r *Repository) GetGithubIssues(token string) ([]*github.Issue, []*github.Issue, error) {
	return hub.GetGithubIssues(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String)
}

func (r *Repository) ListLabels(token string) ([]*github.Label, error) {
	return hub.ListLabels(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String)
}

func (r *Repository) CreateWebhook(token, url string) error {
	return hub.CreateWebhook(token, r.RepositoryModel.Owner.String, r.RepositoryModel.Name.String, r.RepositoryModel.WebhookKey, url)
}
