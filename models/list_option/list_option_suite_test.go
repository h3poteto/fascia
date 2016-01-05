package list_option_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestListOption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ListOption Suite")
}
