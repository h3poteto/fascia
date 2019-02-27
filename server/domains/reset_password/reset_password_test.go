package reset_password_test

import (
	. "github.com/h3poteto/fascia/server/domains/reset_password"
	"github.com/h3poteto/fascia/server/handlers"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyResetPassword{}
}

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
		resetPassword, err = GenerateResetPassword(user.UserEntity.ID, email, dummyInjector())
		if err != nil {
			panic(err)
		}
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			Expect(Authenticate(resetPassword.ID, resetPassword.Token, dummyInjector())).To(BeNil())
		})
	})
})
