package reset_password_test

import (
	. "github.com/h3poteto/fascia/reset_password"

	"database/sql"
	"github.com/h3poteto/fascia/db"
	"github.com/h3poteto/fascia/user"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResetPassword", func() {
	var (
		resetPassword *ResetPasswordStruct
		database      *sql.DB
		password      string
		email         string
	)
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table reset_passwords;")
	})
	JustBeforeEach(func() {
		email = "hoge@example.com"
		password = "hogehoge"
		uid, _ := user.Registration(email, password, password)
		database = db.SharedInstance().Connection
		resetPassword = GenerateResetPassword(uid, email)
		resetPassword.Save()
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
			u, err := ChangeUserPassword(resetPassword.ID, resetPassword.Token, newPassword)
			Expect(err).To(BeNil())
			u, err = user.Login(email, password)
			Expect(err).NotTo(BeNil())
			Expect(u).To(BeNil())
		})
		It("can login with new password", func() {
			u, err := ChangeUserPassword(resetPassword.ID, resetPassword.Token, newPassword)
			Expect(err).To(BeNil())
			u, err = user.Login(email, newPassword)
			Expect(err).To(BeNil())
			Expect(u).NotTo(BeNil())
		})
	})
})
