package reset_password_test

import (
	. "../reset_password"

	"../db"
	"../user"
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResetPassword", func() {
	var (
		resetPassword *ResetPasswordStruct
		table         *sql.DB
		password      string
		email         string
	)
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table reset_passwords;")
		sql.Close()
	})
	JustBeforeEach(func() {
		email = "hoge@example.com"
		password = "hogehoge"
		uid, _ := user.Registration(email, password)
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
		resetPassword = GenerateResetPassword(uid, email)
		resetPassword.Save()
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			Expect(Authenticate(resetPassword.ID, resetPassword.Token)).To(BeTrue())
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
