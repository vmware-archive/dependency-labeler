package kpack_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKpack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kpack Suite")
}
