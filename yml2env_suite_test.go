package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "yml2env Suite")
}
