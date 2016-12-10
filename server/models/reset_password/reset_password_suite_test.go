package reset_password_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestResetPassword(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ResetPassword Suite")
}
