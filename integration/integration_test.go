package integration_test

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with an image reference", func() {
		It("throws an error if scratch image is provided", func() {
			By("executing it")
			inputImage := "scratch"
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", "does-not-matter"}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("'scratch' is a reserved name"))
		})

		It("throws an error if an invalid image sent to docker engine", func() {
			By("executing it")
			inputImage := "swkichtlsmhasd" // random string unlikely for an image ever to exist
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", "does-not-matter"}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("pull access denied for swkichtlsmhasd, repository does not exist or may require 'docker login'"))
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
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", "does-not-matter"}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("invalid reference format"))
		})

		It("returns an image with a tag if the tag flag is provided", func() {
			outputImage, _, _, repoTags := runDeplabAgainstImage("ubuntu:bionic", "--tag", "testtag")

			Expect(repoTags).ToNot(BeEmpty())
			Expect(repoTags).To(ContainElement(ContainSubstring("testtag")))

			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("exits with an error if the tag passed is not valid", func() {

			_, stdErr := runDepLab([]string{"--image", "ubuntu:bionic", "--git", pathToGitRepo, "--tag", "foo:testtag/bar"}, 1)

			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("tag foo:testtag/bar is not valid"))
		})
	})
})
