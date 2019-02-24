package repository_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	. "github.com/h3poteto/fascia/server/domains/entities/repository"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	Describe("Save", func() {
		repositoryID := int64(123456)
		It("should create repository", func() {
			newRepository := New(0, repositoryID, "owner", "repository_name", "test_token")
			err := newRepository.Save()
			Expect(err).To(BeNil())
		})
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			webhookKey := GenerateWebhookKey("repository_name")
			newRepository := New(0, int64(12345), "owner", "repository_name", webhookKey)
			newRepository.Save()
			mac := hmac.New(sha1.New, []byte(webhookKey))
			mac.Write([]byte(""))
			hashedWebhookKey := hex.EncodeToString(mac.Sum(nil))
			Expect(newRepository.Authenticate("sha1="+hashedWebhookKey, []byte(""))).To(BeNil())
		})
	})
})
