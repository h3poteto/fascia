package repository_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"github.com/h3poteto/fascia/db"
	. "github.com/h3poteto/fascia/repository"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		database *sql.DB
	)
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table repositories;")
	})
	JustBeforeEach(func() {
		database = db.SharedInstance().Connection
	})

	Describe("Save", func() {
		repositoryID := int64(123456)
		It("should create repository", func() {
			newRepository := NewRepository(0, repositoryID, "owner", "repository_name", "test_token")
			err := newRepository.Save()
			Expect(err).To(BeNil())
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
			Expect(newRepository.Authenticate("sha1="+hashedWebhookKey, []byte(""))).To(BeNil())
		})
	})
})
