package list_option_test

import (
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	. "github.com/h3poteto/fascia/models/list_option"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListOption", func() {
	BeforeEach(func() {
		seed.ListOptions()
	})
	AfterEach(func() {
		database := db.SharedInstance().Connection
		database.Exec("truncate table list_options;")
	})
	Describe("ListOptionAll", func() {
		It("should list up all list_options", func() {
			options, err := ListOptionAll()
			Expect(err).To(BeNil())
			Expect(options[0].Action).To(Equal("close"))
			Expect(options[1].Action).To(Equal("open"))
		})
	})
})
