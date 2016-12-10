package repository

import (
	"github.com/h3poteto/fascia/server/models/db"

	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Repository interface {
	Save() bool
}

type RepositoryStruct struct {
	ID           int64
	RepositoryID int64
	Owner        sql.NullString
	Name         sql.NullString
	WebhookKey   string
	database     *sql.DB
}

// GenerateWebhookKey create new md5 hash
func GenerateWebhookKey(seed string) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

// New is build new Repository struct
func New(id int64, repositoryID int64, owner string, name string, webhookKey string) *RepositoryStruct {
	if repositoryID <= 0 {
		return nil
	}
	repository := &RepositoryStruct{ID: id, RepositoryID: repositoryID, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}, WebhookKey: webhookKey}
	repository.Initialize()
	return repository
}

// FindRepositoryByRepositoryID is return a Repository struct from repository_id
func FindRepositoryByRepositoryID(repositoryID int64) (*RepositoryStruct, error) {
	database := db.SharedInstance().Connection
	var id int64
	var owner, name, webhookKey string
	err := database.QueryRow("select id, repository_id, owner, name, webhook_key from repositories where repository_id = ?;", repositoryID).Scan(&id, &repositoryID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, errors.Wrap(err, "sql select error")
	}
	return New(id, repositoryID, owner, name, webhookKey), nil
}

func (u *RepositoryStruct) Initialize() {
	u.database = db.SharedInstance().Connection
}

func (u *RepositoryStruct) Save() error {
	result, err := u.database.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", u.RepositoryID, u.Owner, u.Name, u.WebhookKey)
	if err != nil {
		return errors.Wrap(err, "sql execute error")
	}
	u.ID, _ = result.LastInsertId()
	return nil
}

// Authenticate is check token and webhookKey with response
func (u *RepositoryStruct) Authenticate(token string, response []byte) error {
	var webhookKey string
	err := u.database.QueryRow("select webhook_key from repositories where id = ?;", u.ID).Scan(&webhookKey)
	if err != nil {
		return errors.Wrap(err, "sql select error")
	}
	mac := hmac.New(sha1.New, []byte(webhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}
