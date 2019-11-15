package dpkg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDpkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dpkg Suite")
}
