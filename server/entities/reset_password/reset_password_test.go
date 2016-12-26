package reset_password_test

import (
	"database/sql"

	. "github.com/h3poteto/fascia/server/entities/reset_password"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResetPassword", func() {
	var (
		resetPassword *ResetPassword
		database      *sql.DB
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
		database = db.SharedInstance().Connection
		resetPassword, err = GenerateResetPassword(user.UserEntity.UserModel.ID, email)
		err = resetPassword.Save()
		if err != nil {
			panic(err)
		}
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table reset_passwords;")
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			Expect(Authenticate(resetPassword.ResetPasswordModel.ID, resetPassword.ResetPasswordModel.Token)).To(BeNil())
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
