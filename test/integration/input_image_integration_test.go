package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with a distroless base image", func() {
		It("labels the image", func() {
			metadataLabel := runDeplabAgainstImage("gcr.io/distroless/base")

			Expect(metadataLabel.Base).ToNot(BeEmpty())
		})
	})
})
