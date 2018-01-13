package services_test

import (
	. "github.com/h3poteto/fascia/server/services"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	Describe("FindOrCreateUserFromGithub", func() {
		token := os.Getenv("TEST_TOKEN")
		Context("after registration from github", func() {
			user, err := FindOrCreateUserFromGithub(token)
			It("registration suceeded", func() {
				Expect(err).To(BeNil())
				Expect(user).NotTo(BeNil())
				findUser, _ := FindOrCreateUserFromGithub(token)
				Expect(findUser.UserEntity.ID).To(Equal(user.UserEntity.ID))
				Expect(findUser.UserEntity.ID).NotTo(BeZero())
			})
		})
		Context("after regist with email address", func() {
			email := "already_regist@example.com"
			RegistrationUser(email, "hogehoge", "hogehoge")
			user, _ := FindOrCreateUserFromGithub(token)
			It("should update github information", func() {
				Expect(user.UserEntity.OauthToken.Valid).To(BeTrue())
				Expect(user.UserEntity.OauthToken.String).To(Equal(token))
				Expect(user.UserEntity.UUID.Valid).To(BeTrue())
			})
		})
	})
})
