package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
	"strings"
)

var _ = Describe("deplab", func() {
	Context("with an image tar path", func() {
		It("labels the image", func() {
			By("executing it")
			inputTarPath, err := filepath.Abs(filepath.Join("assets", "tiny.tgz"))
			Expect(err).ToNot(HaveOccurred())
			_, _, metadataLabel, _ := runDeplabAgainstTar(inputTarPath)

			gitDependency := filterGitDependency(metadataLabel.Dependencies)

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
			_, stdErr := runDepLab([]string{"--image-tar", inputTarPath, "--git", pathToGitRepo}, 1)
			errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
			Expect(errorOutput).To(ContainSubstring("could not load docker image from tar: open /path/to/image.tar: no such file or directory"))
		})
	})

})
