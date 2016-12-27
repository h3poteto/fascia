package services_test

import (
	"database/sql"
	"os"

	"github.com/h3poteto/fascia/server/models/db"
	. "github.com/h3poteto/fascia/server/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	var (
		database *sql.DB
	)
	BeforeEach(func() {
		database = db.SharedInstance().Connection
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
	})

	Describe("FindOrCreateUserFromGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		Context("after registration from github", func() {
			user, err := FindOrCreateUserFromGithub(token)
			It("registration suceeded", func() {
				Expect(err).To(BeNil())
				Expect(user).NotTo(BeNil())
				findUser, _ := FindOrCreateUserFromGithub(token)
				Expect(findUser.UserEntity.UserModel.ID).To(Equal(user.UserEntity.UserModel.ID))
				Expect(findUser.UserEntity.UserModel.ID).NotTo(BeZero())
			})
		})
		Context("after regist with email address", func() {
			email := "already_regist@example.com"
			RegistrationUser(email, "hogehoge", "hogehoge")
			user, _ := FindOrCreateUserFromGithub(token)
			It("should update github information", func() {
				Expect(user.UserEntity.UserModel.OauthToken.Valid).To(BeTrue())
				Expect(user.UserEntity.UserModel.OauthToken.String).To(Equal(token))
				Expect(user.UserEntity.UserModel.UUID.Valid).To(BeTrue())
			})
		})
	})
})
