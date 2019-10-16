package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/pivotal/deplab/metadata"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pathToBin                 string
	commitHash, pathToGitRepo string
)

func TestDeplab(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var (
			err error
		)

		commitHash, pathToGitRepo = makeFakeGitRepo()

		pathToBin, err = gexec.Build("github.com/pivotal/deplab/cmd/deplab")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
		os.RemoveAll(pathToGitRepo)
		gexec.Kill()
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "deplab Suite")
}

func runDepLab(args []string, expErrCode int) (stdOut *bytes.Reader, stdErr *bytes.Reader) {
	stdOutBuffer := bytes.Buffer{}
	stdErrBuffer := bytes.Buffer{}

	cmd := exec.Command(pathToBin, args...)

	session, err := gexec.Start(cmd, &stdOutBuffer, &stdErrBuffer)
	Expect(err).ToNot(HaveOccurred())
	<-session.Exited

	stdOut = bytes.NewReader(stdOutBuffer.Bytes())
	stdErr = bytes.NewReader(stdErrBuffer.Bytes())

	if os.Getenv("DEBUG") != "" {
		io.Copy(os.Stdout, stdOut)
		io.Copy(os.Stdout, stdErr)
		stdOut.Seek(0, 0)
		stdErr.Seek(0, 0)
	}

	Eventually(session, time.Minute).Should(gexec.Exit(expErrCode))

	return stdOut, stdErr
}

func runDeplabAgainstImage(inputImage string, extraArgs ...string) (metadataLabel metadata.Metadata) {
	f, err := ioutil.TempFile("", "")
	Expect(err).ToNot(HaveOccurred())

	defer os.Remove(f.Name())

	By("executing it")
	args := []string{"--image", inputImage, "--git", pathToGitRepo, "--metadata-file", f.Name()}
	args = append(args, extraArgs...)
	_, _ = runDepLab(args, 0)

	metadataLabel = metadata.Metadata{}
	err = json.NewDecoder(f).Decode(&metadataLabel)
	Expect(err).ToNot(HaveOccurred())

	return metadataLabel
}

func runDeplabAgainstTar(inputTarPath string, extraArgs ...string) (metadataLabel metadata.Metadata) {
	f, err := ioutil.TempFile("", "")
	Expect(err).ToNot(HaveOccurred())

	By("executing it")
	args := []string{"--image-tar", inputTarPath, "--git", pathToGitRepo, "--metadata-file", f.Name()}
	args = append(args, extraArgs...)
	_, _ = runDepLab(args, 0)

	decoder := json.NewDecoder(f)
	metadataLabel = metadata.Metadata{}
	err = decoder.Decode(&metadataLabel)
	Expect(err).ToNot(HaveOccurred())

	return metadataLabel
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

func getContentsOfReader(r io.Reader) []byte {
	contents, err := ioutil.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())

	return contents
}
