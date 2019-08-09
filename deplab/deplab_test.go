package deplab_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"

	. "github.com/pivotal/deplab/deplab"
)

var _ = Describe("Cmd", func() {

	Context("Root", func() {
		Context("GenerateDependencies", func() {
			Context("image with dpkg", func() {
				It("should generate valid json string of dpkg dependencies", func() {
					log.SetOutput(GinkgoWriter)
					dependencies, err := GenerateDependencies("ubuntu:bionic")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(dependencies)).To(Equal(1))

					dependencyMetadata := dependencies[0].Source.Metadata

					dpkgMetadata := dependencyMetadata.(metadata.DebianPackageListSourceMetadata)

					Expect(len(dpkgMetadata.Packages)).To(Equal(89))
				})
			})
			Context("image without dpkg", func() {
				It("should generate valid json string with zero dpkg dependencies", func() {
					log.SetOutput(GinkgoWriter)
					dependencies, err := GenerateDependencies("alpine:latest")
					Expect(err).ToNot(HaveOccurred())

					Expect(len(dependencies)).To(Equal(1))

					dependencyMetadata := dependencies[0].Source.Metadata

					dpkgMetadata := dependencyMetadata.(metadata.DebianPackageListSourceMetadata)

					Expect(len(dpkgMetadata.Packages)).To(Equal(0))
				})
			})
		})
	})
})
