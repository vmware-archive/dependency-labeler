package integration_test

import (
	"context"
	"github.com/pivotal/deplab/metadata"

	"github.com/docker/docker/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab git", func() {
	var outputImage string

	Context("when I supply a git repo as an argument", func() {
		var (
			metadataLabel     metadata.Metadata
			gitDependency     metadata.Dependency
			gitSourceMetadata map[string]interface{}
		)

		BeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage)
			gitDependency = filterGitDependency(metadataLabel.Dependencies)
			gitSourceMetadata = gitDependency.Source.Metadata.(map[string]interface{})
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("adds a gitDependency", func() {
			Expect(gitDependency.Type).ToNot(BeEmpty())

			By("adding the git commit of HEAD to a git dependency")
			Expect(gitDependency.Type).To(Equal("package"))
			Expect(gitDependency.Source.Version["commit"]).To(Equal(commitHash))

			By("adding the git remote to a git dependency")
			Expect(gitSourceMetadata["url"].(string)).To(Equal("https://example.com/example.git"))

			By("adding refs for the current HEAD")
			Expect(len(gitSourceMetadata["refs"].([]interface{}))).To(Equal(1))
			Expect(gitSourceMetadata["refs"].([]interface{})[0].(string)).To(Equal("bar"))

			By("not adding refs that are not the current HEAD")
			Expect(gitSourceMetadata["refs"].([]interface{})[0].(string)).ToNot(Equal("foo"))
		})
	})

	Context("when I supply non-git repo as an argument", func() {
		It("exits with an error message", func() {
			By("executing it")
			inputImage := "ubuntu:bionic"
			_, stdErr := runDepLab([]string{"--image", inputImage, "--git", "/dev/null"}, 1)

			Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("cannot open git repository"))
		})
	})

	Context("when I don't supply a git flag as an argument", func() {
		It("has no git metadata", func() {
			By("executing it")
			inputImage := "ubuntu:bionic"
			_, stdErr := runDepLab([]string{"--image", inputImage}, 1)

			Expect(string(getContentsOfReader(stdErr))).To(ContainSubstring("required flag(s) \"git\" not set"))
		})
	})
})

func filterGitDependency(dependencies []metadata.Dependency) metadata.Dependency {
	for _, dependency := range dependencies {
		if dependency.Source.Type == "git" {
			return dependency
		}
	}
	return metadata.Dependency{}
}
