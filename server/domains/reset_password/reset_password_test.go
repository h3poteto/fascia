package reset_password_test

import (
	. "github.com/h3poteto/fascia/server/domains/entities/reset_password"
	"github.com/h3poteto/fascia/server/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResetPassword", func() {
	var (
		resetPassword *ResetPassword
		password      string
		email         string
	)
	BeforeEach(func() {
		email = "hoge@example.com"
		password = "hogehoge"
		user, err := handlers.RegistrationUser(email, password, password)
		if err != nil {
			panic(err)
		}
		resetPassword, err = GenerateResetPassword(user.UserEntity.ID, email)
		err = resetPassword.Save()
		if err != nil {
			panic(err)
		}
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			Expect(Authenticate(resetPassword.ID, resetPassword.Token)).To(BeNil())
		})
	})

	Describe("ChangeUserPassword", func() {
		var (
			newPassword string
		)
		JustBeforeEach(func() {
			newPassword = "fugafuga"
		})
		It("can not login with old password", func() {
			_, err := resetPassword.ChangeUserPassword(newPassword)
			Expect(err).To(BeNil())
			u, err := handlers.LoginUser(email, password)
			Expect(err).NotTo(BeNil())
			Expect(u).To(BeNil())
		})
		It("can login with new password", func() {
			_, err := resetPassword.ChangeUserPassword(newPassword)
			Expect(err).To(BeNil())
			u, err := handlers.LoginUser(email, newPassword)
			Expect(err).To(BeNil())
			Expect(u).NotTo(BeNil())
		})
	})
})
