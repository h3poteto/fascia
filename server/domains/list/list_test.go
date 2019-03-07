package list_test

import (
	"database/sql"

	. "github.com/h3poteto/fascia/server/domains/list"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyList{
		ProjectID: 2,
		UserID:    1306,
		Option: &dummy.ListOption{
			ID:     1,
			Action: "TODO",
		},
	}
}

var _ = Describe("List", func() {
	Describe("Update", func() {
		It("option does not exist", func() {
			l := New(0, 2, 1306, sql.NullString{String: "title", Valid: true}, sql.NullString{String: "#00ff00", Valid: true}, sql.NullInt64{}, false, dummyInjector())
			err := l.Update(sql.NullString{String: "title2", Valid: true}, sql.NullString{String: "#ff0000", Valid: true}, 0)
			Expect(err).To(BeNil())
			Expect(l.ListOptionID.Valid).To(Equal(false))
		})
		It("option exists", func() {
			l := New(0, 2, 1306, sql.NullString{String: "title", Valid: true}, sql.NullString{String: "#00ff00", Valid: true}, sql.NullInt64{}, false, dummyInjector())
			err := l.Update(sql.NullString{String: "title2", Valid: true}, sql.NullString{String: "#ff0000", Valid: true}, 1)
			Expect(err).To(BeNil())
			Expect(l.ListOptionID.Valid).To(Equal(true))
		})
	})
})
