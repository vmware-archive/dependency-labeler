package rpm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/pkg/rpm"

	"path/filepath"
)



type MockImage struct {}

func (m MockImage) AbsolutePath(string) string {
 	path, err := filepath.Abs("../../test/integration/assets/rpm")

	Expect(err).ToNot(HaveOccurred())
 	return path
}

var _ = Describe("Pkg/Rpm/Provider", func() {
	It("should generate list of dependencies", func() {
		md, err := rpm.BuildRPMMetadata(MockImage{})

		Expect(err).ToNot(HaveOccurred())
		Expect(md.Source.Metadata.(metadata.RpmPackageListSourceMetadata).Packages).To(HaveLen(34))
	})
})
