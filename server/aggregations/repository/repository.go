package repository

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/h3poteto/fascia/lib/modules/hub"
	"github.com/h3poteto/fascia/server/models/repository"
)

type Repository struct {
	RepositoryModel *repository.RepositoryStruct
}

func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *Repository {
	return &Repository{
		RepositoryModel: repository.New(id, repositoryID, owner, name, webhookKey),
	}
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
