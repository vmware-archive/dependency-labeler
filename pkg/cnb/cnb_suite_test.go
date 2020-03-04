package cnb_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCnb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cnb Suite")
}
