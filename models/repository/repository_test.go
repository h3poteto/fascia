package repository_test

import (
	"../db"
	. "../repository"
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		table *sql.DB
	)
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table repositories;")
		sql.Close()
	})
	JustBeforeEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
	})

	Describe("Save", func() {
		repositoryID := int64(123456)
		It("should create repository", func() {
			newRepository := NewRepository(0, repositoryID, "owner", "repository_name", "test_token")
			result := newRepository.Save()
			Expect(result).To(BeTrue())
		})
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			webhookKey := GenerateWebhookKey("repository_name")
			newRepository := NewRepository(0, int64(12345), "owner", "repository_name", webhookKey)
			newRepository.Save()
			mac := hmac.New(sha1.New, []byte(webhookKey))
			mac.Write([]byte(""))
			hashedWebhookKey := hex.EncodeToString(mac.Sum(nil))
			Expect(newRepository.Authenticate("sha1="+hashedWebhookKey, []byte(""))).To(BeTrue())
		})
	})
})
