package repository_test

import (
	"../db"
	. "../repository"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		table *sql.DB
	)
	AfterEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		sql := database.Init()
		sql.Exec("truncate table users;")
		sql.Exec("truncate table repositories;")
		sql.Close()
	})
	JustBeforeEach(func() {
		mydb := &db.Database{}
		var database db.DB = mydb
		table = database.Init()
	})

	Describe("Save", func() {
		repositoryId := int64(123456)
		It("should create repository", func() {
			newRepository := NewRepository(0, repositoryId, "owner", "repository_name")
			result := newRepository.Save()
			Expect(result).To(BeTrue())
		})
	})
})
