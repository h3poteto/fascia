package project_test

import (
	"database/sql"

	. "github.com/h3poteto/fascia/server/domains/project"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyProject{
		UserID: 1306,
	}
}

var _ = Describe("Project", func() {
	Describe("CheckOwner", func() {
		It("owner", func() {
			p := New(0, 1306, "title", "description", sql.NullInt64{}, true, true, dummyInjector())
			owner := p.CheckOwner(1306)
			Expect(owner).To(Equal(true))
		})
		It("not owner", func() {
			p := New(0, 1306, "title", "description", sql.NullInt64{}, true, true, dummyInjector())
			owner := p.CheckOwner(1)
			Expect(owner).To(Equal(false))
		})
	})
})
