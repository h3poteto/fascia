package task_test

import (
	"database/sql"

	. "github.com/h3poteto/fascia/server/domains/task"
	dummy "github.com/h3poteto/fascia/server/test/helpers/repositories"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func dummyInjector() Repository {
	return &dummy.DummyTask{
		ListID:    12,
		ProjectID: 2,
		UserID:    1306,
	}
}

var _ = Describe("Task", func() {
	Describe("ChangeList", func() {
		It("reorder", func() {
			t := New(1, 12, 2, 1306, sql.NullInt64{}, "title", "description", false, sql.NullString{}, dummyInjector())
			prev := int64(4)
			isReorder, err := t.ChangeList(12, &prev)
			Expect(err).To(BeNil())
			Expect(isReorder).To(Equal(true))
		})

		It("not reorder", func() {
			t := New(1, 12, 2, 1306, sql.NullInt64{}, "title", "description", false, sql.NullString{}, dummyInjector())
			prev := int64(4)
			isReorder, err := t.ChangeList(13, &prev)
			Expect(err).To(BeNil())
			Expect(isReorder).To(Equal(false))
		})
	})

	Describe("Delete", func() {
		It("related issue", func() {
			t := New(1, 12, 2, 1306, sql.NullInt64{Int64: 1, Valid: true}, "title", "description", false, sql.NullString{}, dummyInjector())
			err := t.Delete()
			Expect(err).NotTo(BeNil())
		})
		It("not related issue", func() {
			t := New(1, 12, 2, 1306, sql.NullInt64{}, "title", "description", false, sql.NullString{}, dummyInjector())
			err := t.Delete()
			Expect(err).To(BeNil())
		})
	})
})
