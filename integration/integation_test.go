package integration_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	It("throws an error if scratch image is provided", func() {
		By("executing it")
		inputImage := "scratch"
		_, stdErrBuffer := runDepLab([]string{"--image", inputImage}, 1)
		errorOutput := strings.TrimSpace(stdErrBuffer.String())
		Expect(errorOutput).To(ContainSubstring("'scratch' is a reserved name."))
	})

	It("throws an error if an invalid image sent to docker engine", func() {
		By("executing it")
		inputImage := "swkichtlsmhasd" // random string unlikely for an image ever to exist
		_, stdErrBuffer := runDepLab([]string{"--image", inputImage}, 1)

		errorOutput := strings.TrimSpace(stdErrBuffer.String())
		Expect(errorOutput).To(ContainSubstring("pull access denied for swkichtlsmhasd, repository does not exist or may require 'docker login'"))
	})

	It("throws an error if missing parameters", func() {
		By("executing it")
		_, stdErrBuffer := runDepLab([]string{}, 1)

		errorOutput := strings.TrimSpace(stdErrBuffer.String())
		Expect(errorOutput).To(ContainSubstring("required flag(s) \"image\" not set"))
	})

	It("throws an error if invalid characters are in image name", func() {
		By("executing it")
		inputImage := "£$Invalid_image_name$£"
		_, stdErrBuffer := runDepLab([]string{"--image", inputImage}, 1)

		errorOutput := strings.TrimSpace(stdErrBuffer.String())
		Expect(errorOutput).To(ContainSubstring("invalid reference format"))
	})

})
