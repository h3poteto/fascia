package user_test

import (
	. "github.com/h3poteto/fascia/server/domains/user"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyUser{}
}

var _ = Describe("User", func() {
	Describe("Registration", func() {
		email := "registration@example.com"
		password := "samplepassword"
		It("can regist", func() {
			user, err := Registration(email, password, password, dummyInjector())
			Expect(err).To(BeNil())
			Expect(user.ID).To(Equal(int64(1)))
		})
	})

	Describe("Login", func() {
		email := "login@example.com"
		password := "samplepassword"
		Context("when send correctly login information", func() {
			It("can login", func() {
				currentUser, err := Login(email, password, dummyInjector())
				Expect(err).To(BeNil())
				Expect(currentUser.Email).To(Equal(email))
			})
		})
	})
})
