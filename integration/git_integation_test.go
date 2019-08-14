package integration_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4/config"

	"gopkg.in/src-d/go-git.v4"

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
			commitHash        string
			pathToGitRepo     string
			gitDependency     metadata.Dependency
			gitSourceMetadata map[string]interface{}
		)

		BeforeEach(func() {
			inputImage := "ubuntu:bionic"
			commitHash, pathToGitRepo = makeFakeGitRepo()
			outputImage, _, metadataLabel = runDeplabAgainstImage(inputImage, "--git", pathToGitRepo)
			gitDependency = filterGitDependency(metadataLabel.Dependencies)
			gitSourceMetadata = gitDependency.Source.Metadata.(map[string]interface{})
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
			os.RemoveAll(pathToGitRepo)
		})

		It("adds a gitDependency", func() {
			Expect(gitDependency.Type).ToNot(BeEmpty())

			By("adding the git commit of HEAD to a git dependency")
			Expect(gitDependency.Type).To(Equal("package"))
			Expect(gitDependency.Source.Version["commit"]).To(Equal(commitHash))

			By("adding the git remote to a git dependency")
			Expect(gitSourceMetadata["uri"].(string)).To(Equal("https://example.com/example.git"))

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
			_, stdErrBuffer := runDepLab([]string{"--image", inputImage, "--git", "/dev/null"}, 1)

			Expect(stdErrBuffer.String()).To(ContainSubstring("cannot open git repository"))
		})
	})

	Context("when I don't supply a git flag as an argument", func() {
		var (
			gitDependency metadata.Dependency
			metadataLabel metadata.Metadata
		)

		BeforeEach(func() {
			inputImage := "ubuntu:bionic"
			outputImage, _, metadataLabel = runDeplabAgainstImage(inputImage)
			gitDependency = filterGitDependency(metadataLabel.Dependencies)
		})

		AfterEach(func() {
			_, err := dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
		})

		It("has no git metadata", func() {
			Expect(gitDependency.Type).To(BeEmpty())
		})
	})
})

func makeFakeGitRepo() (string, string) {
	path, err := ioutil.TempDir("", "deplab-integration")
	Expect(err).ToNot(HaveOccurred())

	repo, err := git.PlainInit(path, false)
	Expect(err).ToNot(HaveOccurred())

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://example.com/example.git"},
	})
	Expect(err).ToNot(HaveOccurred())

	testFilePath := filepath.Join(path, "test")
	data := []byte("TestFile\n")
	err = ioutil.WriteFile(testFilePath, data, 0644)
	Expect(err).ToNot(HaveOccurred())

	w, err := repo.Worktree()
	Expect(err).ToNot(HaveOccurred())

	err = w.AddGlob("*")
	Expect(err).ToNot(HaveOccurred())

	ch, err := w.Commit("Test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Pivotal Example",
			Email: "example@pivotal.io",
			When:  time.Now(),
		},
	})
	Expect(err).ToNot(HaveOccurred())

	repo.CreateTag("foo", ch, nil)

	ch, err = w.Commit("Second test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Pivotal Example",
			Email: "example@pivotal.io",
			When:  time.Now(),
		},
	})
	Expect(err).ToNot(HaveOccurred())

	repo.CreateTag("bar", ch, nil)

	return ch.String(), path
}

func filterGitDependency(dependencies []metadata.Dependency) metadata.Dependency {
	for _, dependency := range dependencies {
		if dependency.Source.Type == "git" {
			return dependency
		}
	}
	return metadata.Dependency{} //should never be reached
}
