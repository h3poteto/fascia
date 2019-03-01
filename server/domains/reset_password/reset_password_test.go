package reset_password_test

import (
	. "github.com/h3poteto/fascia/server/domains/reset_password"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyResetPassword{}
}

var _ = Describe("ResetPassword", func() {
	BeforeEach(func() {
		_, err := GenerateResetPassword(1, "hoge@example.com", dummyInjector())
		if err != nil {
			panic(err)
		}
	})

	Describe("Authenticate", func() {
		It("should authenticate", func() {
			Expect(Authenticate(1, "test token", dummyInjector())).To(BeNil())
		})
	})
})
