package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFascia(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fascia Suite")
}
