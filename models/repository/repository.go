package repository

import (
	"../../modules/logging"
	"../db"
	"crypto/md5"
	"database/sql"
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

func GenerateWebhookKey(seed string) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, seed)
	token := fmt.Sprintf("%x", h.Sum(nil))

	return token
}

func NewRepository(id int64, repositoryID int64, owner string, name string, webhookKey string) *RepositoryStruct {
	if repositoryID <= 0 {
		return nil
	}
	repository := &RepositoryStruct{ID: id, RepositoryID: repositoryID, Owner: sql.NullString{String: owner, Valid: true}, Name: sql.NullString{String: name, Valid: true}, WebhookKey: webhookKey}
	repository.Initialize()
	return repository
}

func (u *RepositoryStruct) Initialize() {
	objectDB := &db.Database{}
	var interfaceDB db.DB = objectDB
	u.database = interfaceDB
}

func (u *RepositoryStruct) Save() bool {
	table := u.database.Init()
	defer table.Close()

	result, err := table.Exec("insert into repositories (repository_id, owner, name, webhook_key, created_at) values (?, ?, ?, ?, now());", u.RepositoryID, u.Owner, u.Name, u.WebhookKey)
	if err != nil {
		logging.SharedInstance().MethodInfo("Repository", "Save", true).Errorf("repository save failed: %v", err)
		return false
	}
	u.ID, _ = result.LastInsertId()
	return true
}

func (u *RepositoryStruct) Authenticate() bool {
	table := u.database.Init()
	defer table.Close()

	var id int64
	err := table.QueryRow("select id from repositories where id = ? and webhook_key = ?;", u.ID, u.WebhookKey).Scan(&id)
	if err != nil {
		logging.SharedInstance().MethodInfo("Repository", "Authenticate").Infof("cannot authenticate to repository webhook_key: %v", err)
		return false
	}
	return true
}
