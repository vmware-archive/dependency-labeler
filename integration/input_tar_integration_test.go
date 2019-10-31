package integration_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with an image tar path", func() {
		It("labels the image", func() {
			By("executing it")
			inputTarPath, err := filepath.Abs(filepath.Join("assets", "tiny.tgz"))
			Expect(err).ToNot(HaveOccurred())
			metadataLabel := runDeplabAgainstTar(inputTarPath)

			gitDependencies := selectGitDependencies(metadataLabel.Dependencies)
			gitDependency := gitDependencies[0]

			Expect(gitDependency.Type).ToNot(BeEmpty())

			By("adding the git commit of HEAD to a git dependency")
			Expect(gitDependency.Type).To(Equal("package"))
			Expect(gitDependency.Source.Version["commit"]).To(Equal(commitHash))
		})
	})

	Context("with an invalid image tar path", func() {
		It("exits with an error", func() {
			By("executing it")
			inputTarPath := "/path/to/image.tar"
			_, stdErr := runDepLab([]string{
				"--image-tar", inputTarPath,
				"--git", pathToGitRepo,
				"--metadata-file", "doesnotmatter7",
			}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(
				SatisfyAll(
					ContainSubstring("/path/to/image.tar"),
					ContainSubstring("could not load image"),
				))
		})
	})

	Context("with an image tar with a eu.gcr.io tag", func() {
		It("labels the image", func() {
			By("executing it")
			inputTarPath, err := filepath.Abs(filepath.Join("assets", "tiny-with-eu.gcr.io-tag.tar"))
			Expect(err).ToNot(HaveOccurred())
			metadataLabel := runDeplabAgainstTar(inputTarPath)

			gitDependencies := selectGitDependencies(metadataLabel.Dependencies)
			gitDependency := gitDependencies[0]

			Expect(gitDependency.Type).ToNot(BeEmpty())

			By("adding the git commit of HEAD to a git dependency")
			Expect(gitDependency.Type).To(Equal("package"))
			Expect(gitDependency.Source.Version["commit"]).To(Equal(commitHash))
		})
	})

})
