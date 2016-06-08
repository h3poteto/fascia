package repository

import (
	"../db"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
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
	database     db.DB
}

// GenerateWebhookKey is create new md5 hash
func GenerateWebhookKey(seed string) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

// NewRepository is build new Repository struct
func NewRepository(id int64, repositoryID int64, owner string, name string, webhookKey string) *RepositoryStruct {
	if repositoryID <= 0 {
		return nil
	}
	repository := &RepositoryStruct{ID: id, RepositoryID: repositoryID, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}, WebhookKey: webhookKey}
	repository.Initialize()
	return repository
}

// FindRepositoryByRepositoryID is return a Repository struct from repository_id
func FindRepositoryByRepositoryID(repositoryID int64) (*RepositoryStruct, error) {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	table := interfaceDB.Init()
	defer table.Close()

	var id int64
	var owner, name, webhookKey string
	err := table.QueryRow("select id, repository_id, owner, name, webhook_key from repositories where repository_id = ?;", repositoryID).Scan(&id, &repositoryID, &owner, &name, &webhookKey)
	if err != nil {
		return nil, err
	}
	return NewRepository(id, repositoryID, owner, name, webhookKey), nil
}

func (u *RepositoryStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

// Validation check record required column
func (u *RepositoryStruct) Validation(tx *sql.Tx) bool {
	if u.RepositoryID == 0 {
		return false
	}
	return true
}

func (u *RepositoryStruct) Save() error {
	table := u.database.Init()
	defer table.Close()
	tx, err := table.Begin()
	if err != nil {
		panic(err)
	}

	if !u.Validation(tx) {
		return errors.New("validation error")
	}

	result, err := tx.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", u.RepositoryID, u.Owner, u.Name, u.WebhookKey)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
	u.ID, _ = result.LastInsertId()
	return nil
}

// Authenticate is check token and webhookKey with response
func (u *RepositoryStruct) Authenticate(token string, response []byte) error {
	table := u.database.Init()
	defer table.Close()

	var webhookKey string
	err := table.QueryRow("select webhook_key from repositories where id = ?;", u.ID).Scan(&webhookKey)
	if err != nil {
		return err
	}
	mac := hmac.New(sha1.New, []byte(webhookKey))
	mac.Write(response)
	hashedToken := hex.EncodeToString(mac.Sum(nil))
	if token != ("sha1=" + hashedToken) {
		return errors.New("token is not equal webhookKey")
	}
	return nil
}
