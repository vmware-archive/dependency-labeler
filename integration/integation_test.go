package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4/config"

	"gopkg.in/src-d/go-git.v4"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/docker/docker/api/types"

	docker "github.com/docker/docker/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var outputImage string

func runDepLab(args []string, expErrCode int) (stdOutBuffer bytes.Buffer, stdErrBuffer bytes.Buffer) {
	stdOutBuffer = bytes.Buffer{}
	stdErrBuffer = bytes.Buffer{}

	cmd := exec.Command(pathToBin, args...)

	session, err := gexec.Start(cmd, &stdOutBuffer, &stdErrBuffer)
	Expect(err).ToNot(HaveOccurred())
	<-session.Exited

	Eventually(session, time.Minute).Should(gexec.Exit(expErrCode))

	return stdOutBuffer, stdErrBuffer
}

var _ = Describe("deplab", func() {
	dockerCli, err := docker.NewClientWithOpts(docker.WithVersion("1.39"), docker.FromEnv)
	if err != nil {
		panic(err)
	}

	It("labels an image and returns the sha of the labelled image with a dpkg list", func() {
		By("executing it")
		inputImage := "ubuntu:bionic"
		stdOutBuffer, _ := runDepLab([]string{"--image", inputImage}, 0)

		By("checking if it returns an image sha")
		outputImage = strings.TrimSpace(stdOutBuffer.String())
		Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

		By("checking if the label exists")
		inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
		Expect(err).ToNot(HaveOccurred())

		labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
		Expect(labelValue).ToNot(BeEmpty())

		By("checking if the dpkg dependencies exists")
		result := metadata.Metadata{}
		err = json.Unmarshal([]byte(labelValue), &result)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(result.Dependencies)).To(Equal(1))

		dependencyMetadata := result.Dependencies[0].Source.Metadata
		dpkgMetadata := dependencyMetadata.(map[string]interface{})
		Expect(len(dpkgMetadata["pkg"].([]interface{}))).To(Equal(89))

		By("checking that the input image is parent of the output image")
		inspectInput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), inputImage)
		Expect(err).ToNot(HaveOccurred())

		Expect(inspectOutput.Parent).To(Equal(inspectInput.ID))

		_, err = dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with an image without dpkg", func() {
		It("does not return a dpkg list", func() {
			By("executing it")
			inputImage := "alpine:latest"
			stdOutBuffer, _ := runDepLab([]string{"--image", inputImage}, 0)

			By("checking if it returns an image sha")
			outputImage = strings.TrimSpace(stdOutBuffer.String())
			Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

			By("checking if the label exists")
			inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
			Expect(err).ToNot(HaveOccurred())

			labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
			Expect(labelValue).ToNot(BeEmpty())

			By("checking if the dpkg dependencies exists")
			result := metadata.Metadata{}
			err = json.Unmarshal([]byte(labelValue), &result)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(result.Dependencies)).To(Equal(0))
		})
	})
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
		Expect(errorOutput).To(ContainSubstring("invalid image name"))
	})

	Context("when I supply a git repo as an argument", func() {
		It("adds git metadata to the metadata label", func() {
			By("executing it")
			inputImage := "ubuntu:bionic"
			commitHash, pathToGitRepo := makeFakeGitRepo()
			stdOutBuffer, _ := runDepLab([]string{"--image", inputImage, "--git", pathToGitRepo}, 0)

			By("checking if it returns an image sha")
			outputImage = strings.TrimSpace(stdOutBuffer.String())
			Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

			By("checking if the label exists")
			inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
			Expect(err).ToNot(HaveOccurred())

			labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
			Expect(labelValue).ToNot(BeEmpty())

			By("checking if two dependencies exists")
			result := metadata.Metadata{}
			err = json.Unmarshal([]byte(labelValue), &result)
			Expect(err).ToNot(HaveOccurred())

			By("checking if a git dependency exists")
			gitDependency := filterGitDependency(result.Dependencies)
			Expect(gitDependency.Type).To(Equal("package"))
			Expect(gitDependency.Source.Version["commit"]).To(Equal(commitHash))

			By("checking a remote is specified")
			gitSourceMetadata := gitDependency.Source.Metadata.(map[string]interface{})
			Expect(gitSourceMetadata["uri"].(string)).To(Equal("https://example.com/example.git"))

			By("checking only 1 tag is included in refs")
			Expect(len(gitSourceMetadata["refs"].([]interface{}))).To(Equal(1))
			Expect(gitSourceMetadata["refs"].([]interface{})[0].(string)).To(Equal("bar"))

			_, err = dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
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

	It("labels an image and returns the sha of the labelled image with a base image", func() {
		By("executing it")
		inputImage := "ubuntu:bionic"
		stdOutBuffer, _ := runDepLab([]string{"--image", inputImage}, 0)

		By("checking if it returns an image sha")
		outputImage = strings.TrimSpace(stdOutBuffer.String())
		Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

		By("checking if the label exists")
		inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
		Expect(err).ToNot(HaveOccurred())

		labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
		Expect(labelValue).ToNot(BeEmpty())

		By("checking if the base metadata exists")
		result := metadata.Metadata{}
		err = json.Unmarshal([]byte(labelValue), &result)

		Expect(result.Base.Name).To(Equal("Ubuntu"))
		Expect(result.Base.VersionCodename).To(Equal("bionic"))
		Expect(result.Base.VersionID).To(Equal("18.04"))

		_, err = dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with an image that dosen't have an os-release", func() {
		It("labels an image and returns the sha of the labelled image with null instead of base metadata", func() {
			By("executing it")
			inputImage := "pivotalnavcon/noosrelease"
			stdOutBuffer, stdErrBuffer := runDepLab([]string{"--image", inputImage}, 0)

			By("checking if it returns an image sha")
			outputImage = strings.TrimSpace(stdOutBuffer.String())
			Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

			Expect(stdErrBuffer.String()).To(ContainSubstring("WARNING: error getting OS info"))

			By("checking if the label exists")
			inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
			Expect(err).ToNot(HaveOccurred())

			labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
			Expect(labelValue).ToNot(BeEmpty())

			By("checking if the base metadata exists")
			result := metadata.Metadata{}
			err = json.Unmarshal([]byte(labelValue), &result)

			Expect(result.Base).To(BeNil())

			_, err = dockerCli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
			Expect(err).ToNot(HaveOccurred())
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
	Fail("Could not find a git dependency")
	return metadata.Dependency{} //should never be reached
}
