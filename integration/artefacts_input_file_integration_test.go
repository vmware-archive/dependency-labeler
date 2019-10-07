package integration_test

import (
	"context"
	"github.com/docker/docker/api/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/metadata"
	"path/filepath"
	"strings"
)

var _ = Describe("deplab artefacts", func(){

	Context("when I supply an artefacts file as an argument", func() {
		var (
			metadataLabel       metadata.Metadata
			additionalArguments []string
			outputImage         string
		)

		JustBeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel, _ = runDeplabAgainstImage(inputImage, additionalArguments...)
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when I supply an artefacts file with only one blob", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-single-blob.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds a blob dependency", func() {
				blobDependencies := selectBlobDependencies(metadataLabel.Dependencies)
				Expect(len(blobDependencies)).To(Equal(1))
				blobDependency := blobDependencies[0]
				Expect(blobDependency.Source.Metadata).NotTo(BeNil())
				blobSourceMetadata := blobDependency.Source.Metadata.(map[string]interface{})
				Expect(blobDependency.Type).ToNot(BeEmpty())

				By("adding the blob url to the blob dependency")
				Expect(blobDependency.Type).To(Equal("package"))
				Expect(blobDependency.Source.Type).To(Equal("blob"))
				Expect(blobSourceMetadata["url"]).To(Equal("http://archive.ubuntu.com/ubuntu/pool/main/c/ca-certificates/ca-certificates_20180409.tar.xz"))
			})
		})

		Context("when I supply an artefacts file with multiple blobs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-multiple-blobs.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds multiple blobDependency entries", func() {

			})
		})

		Context("when I supply an artefacts file with no blobs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-empty-blobs.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds zero blobDependency entries", func() {
				blobDependencies := selectBlobDependencies(metadataLabel.Dependencies)
				Expect(len(blobDependencies)).To(Equal(0))
			})
		})

		Context("when I supply an artefacts file with only one vcs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-single-vcs.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds a git dependency", func() {
				gitDependencies := selectGitDependencies(metadataLabel.Dependencies)
				Expect(len(gitDependencies)).To(Equal(2))
				vcsGitDependencies := selectVcsGitDependencies(gitDependencies)
				Expect(len(vcsGitDependencies)).To(Equal(1))

				gitDependency := vcsGitDependencies[0]
				gitSourceMetadata := gitDependency.Source.Metadata.(map[string]interface{})
				Expect(gitDependency.Type).ToNot(BeEmpty())

				By("adding the git commit of HEAD to a git dependency")
				Expect(gitDependency.Type).To(Equal("package"))
				Expect(gitDependency.Source.Version["commit"]).To(Equal("abc123"))

				By("adding the git remote to a git dependency")
				Expect(gitSourceMetadata["url"].(string)).To(Equal("git@github.com:pivotal/deplab.git"))

				By("adding refs for the current HEAD")
				Expect(len(gitSourceMetadata["refs"].([]interface{}))).To(Equal(1))
				Expect(gitSourceMetadata["refs"].([]interface{})[0].(string)).To(Equal("v0.44.0"))
			})
		})

		Context("when I supply an artefacts file with both multiple vcs and multiple blobs", func() {
			BeforeEach(func() {
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-multiple-blobs-multiple-vcs.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath}
			})

			It("adds git dependencies and blobs", func(){
				gitDependencies := selectGitDependencies(metadataLabel.Dependencies)
				Expect(len(gitDependencies)).To(Equal(3))
				vcsGitDependencies := selectVcsGitDependencies(gitDependencies)
				Expect(len(vcsGitDependencies)).To(Equal(2))
			})
		})

		Context("when I supply multiple artefacts files", func() {
			BeforeEach(func() {
				inputArtefactsPath1, err := filepath.Abs(filepath.Join("assets", "artefacts-multiple-blobs.yml"))
				Expect(err).ToNot(HaveOccurred())
				inputArtefactsPath2, err := filepath.Abs(filepath.Join("assets", "artefacts-single-blob.yml"))
				Expect(err).ToNot(HaveOccurred())
				additionalArguments = []string{"--artefacts-file", inputArtefactsPath1, "--artefacts-file", inputArtefactsPath2}
			})

			It("adds multiple blobDependency entries", func() {
				i := 0
				for _, dep := range metadataLabel.Dependencies {
					if dep.Source.Type == "blob" {
						i++
					}
				}

				Expect(i).To(Equal(3))
			})
		})

		Context("when I supply erroneous paths as artefacts file", func(){
			It("exits with an error", func() {
				By("executing it")
				inputTarPath, err := filepath.Abs(filepath.Join("assets", "tiny.tgz"))
				Expect(err).ToNot(HaveOccurred())
				_, stdErr := runDepLab([]string{"--artefacts-file", "erroneous_path.yml", "--image-tar", inputTarPath, "--git", pathToGitRepo}, 1)
				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
				Expect(errorOutput).To(ContainSubstring("could not parse artefact file: erroneous_path.yml"))
			})
		})

		Context("when I supply empty file as artefacts file", func(){
			It("exits with an error", func() {
				By("executing it")
				inputTarPath, err := filepath.Abs(filepath.Join("assets", "tiny.tgz"))
				Expect(err).ToNot(HaveOccurred())
				inputArtefactsPath, err := filepath.Abs(filepath.Join("assets", "artefacts-empty.yml"))
				Expect(err).ToNot(HaveOccurred())
				_, stdErr := runDepLab([]string{"--artefacts-file", inputArtefactsPath, "--image-tar", inputTarPath, "--git", pathToGitRepo}, 1)
				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))
				Expect(errorOutput).To(ContainSubstring("could not parse artefact file"))
			})
		})
	})


})

func selectVcsGitDependencies(dependencies []metadata.Dependency) []metadata.Dependency {
	var gitDependencies []metadata.Dependency
	for _, dependency := range dependencies {
		Expect(dependency.Source.Metadata).To(Not(BeNil()))
		gitSourceMetadata := dependency.Source.Metadata.(map[string]interface{})
		if dependency.Source.Type == "git" && gitSourceMetadata["url"].(string) != "https://example.com/example.git" {
			gitDependencies = append(gitDependencies, dependency)
		}
	}
	return gitDependencies
}
