package repository_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	AfterEach(func() {
		err := seed.TruncateAll()
		if err != nil {
			panic(err)
		}
	})
	RunSpecs(t, "Repository Suite")
}
