package integration_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with an image reference", func() {
		It("throws an error if scratch image is provided", func() {
			By("executing it")
			inputImage := "scratch"
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(
				SatisfyAll(
					ContainSubstring("could not load image"),
					ContainSubstring("scratch"),
				))
		})

		It("throws an error if an invalid image sent to docker engine", func() {
			By("executing it")
			inputImage := "swkichtlsmhasd" // random string unlikely for an image ever to exist
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(
				SatisfyAll(
					ContainSubstring("could not load image"),
					ContainSubstring("swkichtlsmhasd"),
				))
		})

		It("exits with an error if neither image or image-tar flags are set", func() {
			_, stdErr := runDepLab([]string{"--git", "does-not-matter"}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("ERROR: requires one of --image or --image-tar"))
		})

		It("exits with an error if both image and image-tar flags are set", func() {
			_, stdErr := runDepLab([]string{"--image", "foo", "--image-tar", "path/to/image.tar", "--git", "does-not-matter"}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("ERROR: cannot accept both --image and --image-tar"))
		})

		It("throws an error if invalid characters are in image name", func() {
			By("executing it")
			inputImage := "£$Invalid_image_name$£"
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("could not parse reference"))
		})

		It("exits with an error if additional-source-url is not valid", func() {
			_, stdErr := runDepLab([]string{
				"--image", "ubuntu:bionic",
				"--git", pathToGitRepo,
				"--additional-source-url", "/foo/bar",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(SatisfyAll(
				ContainSubstring("/foo/bar"),
				ContainSubstring("error validating additional source url")))
		})

		It("exits with an error if additional-source-url is not reachable ", func() {
			_, stdErr := runDepLab([]string{
				"--image", "ubuntu:bionic",
				"--git", pathToGitRepo,
				"--additional-source-url", "https://package.some.invalid/cool-package",
			}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(SatisfyAll(
				ContainSubstring("https://package.some.invalid/cool-package"),
				ContainSubstring("error validating additional source url")))
		})
	})
})
