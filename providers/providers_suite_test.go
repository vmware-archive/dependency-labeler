package providers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProviders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Providers Suite")
}
