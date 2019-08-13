package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	docker "github.com/docker/docker/client"
)

var (
	pathToBin string
	dockerCli *docker.Client
)

func TestDeplab(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var (
			err error
		)

		dockerCli, err = docker.NewClientWithOpts(docker.WithVersion("1.39"), docker.FromEnv)
		if err != nil {
			panic(err)
		}

		pathToBin, err = gexec.Build("github.com/pivotal/deplab")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterSuite(func() {
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

func runDeplabAgainstImage(inputImage string, extraArgs ...string) (string, string, metadata.Metadata) {
	By("executing it")
	args := []string{"--image", inputImage}
	args = append(args, extraArgs...)
	stdOutBuffer, _ := runDepLab(args, 0)

	By("checking if it returns an image sha")
	outputImage := strings.TrimSpace(stdOutBuffer.String())
	Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

	By("checking if the label exists")
	inspectOutput, _, err := dockerCli.ImageInspectWithRaw(context.TODO(), outputImage)
	Expect(err).ToNot(HaveOccurred())

	metadataLabelString := inspectOutput.Config.Labels["io.pivotal.metadata"]
	metadataLabel := metadata.Metadata{}
	err = json.Unmarshal([]byte(metadataLabelString), &metadataLabel)
	Expect(err).ToNot(HaveOccurred())

	return outputImage, metadataLabelString, metadataLabel
}
