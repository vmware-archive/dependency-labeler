package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/pivotal/deplab/metadata"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	docker "github.com/docker/docker/client"
)

var (
	pathToBin                 string
	dockerCli                 *docker.Client
	commitHash, pathToGitRepo string
)

func TestDeplab(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var (
			err error
		)

		commitHash, pathToGitRepo = makeFakeGitRepo()

		dockerCli, err = docker.NewClientWithOpts(docker.WithVersion("1.39"), docker.FromEnv)
		if err != nil {
			panic(err)
		}

		pathToBin, err = gexec.Build("github.com/pivotal/deplab")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
		os.RemoveAll(pathToGitRepo)
		gexec.Kill()
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Deplab Suite")
}

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

func runDeplabAgainstImage(inputImage string, extraArgs ...string) (outputImage string, metadataLabelString string, metadataLabel metadata.Metadata) {
	By("executing it")
	args := []string{"--image", inputImage, "--git", pathToGitRepo}
	args = append(args, extraArgs...)
	stdOutBuffer, _ := runDepLab(args, 0)

	return parseOutputAndValidate(stdOutBuffer)
}

func parseOutputAndValidate(stdOutBuffer bytes.Buffer) (outputImage string, metadataLabelString string, metadataLabel metadata.Metadata) {
	By("checking if it returns an image sha")
	outputImage = strings.TrimSpace(stdOutBuffer.String())
	Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

	By("checking if the label exists")
	inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
	Expect(err).ToNot(HaveOccurred())

	metadataLabelString = inspectOutput.Config.Labels["io.pivotal.metadata"]
	metadataLabel = metadata.Metadata{}
	err = json.Unmarshal([]byte(metadataLabelString), &metadataLabel)
	Expect(err).ToNot(HaveOccurred())

	return outputImage, metadataLabelString, metadataLabel
}

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
