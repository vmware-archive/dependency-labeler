package rpm

import (
"testing"

. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
)

func TestRPM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RPM Suite")
}
