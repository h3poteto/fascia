package list_option_test

import (
	"../db"
	. "../list_option"

	seed "../../db/seed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListOption", func() {
	BeforeEach(func() {
		seed.ListOptions()
	})
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table list_options;")
		sql.Close()
	})
	Describe("ListOptionAll", func() {
		It("should list up all list_options", func() {
			options := ListOptionAll()
			Expect(options[0].Action).To(Equal("close"))
			Expect(options[1].Action).To(Equal("open"))
		})
	})
})
