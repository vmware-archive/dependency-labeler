package main_test

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var outputImage string

var _ = Describe("deplab", func() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"), client.FromEnv)
	if err != nil {
		panic(err)
	}

	It("labels an image and returns the new sha", func() {
		stdOutBuffer := bytes.Buffer{}

		By("executing it")
		inputImage := "alpine"
		cmd := exec.Command(pathToBin, inputImage)
		session, err := gexec.Start(cmd, &stdOutBuffer, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, time.Second*5).Should(gexec.Exit(0))
		<-session.Exited

		By("checking if it returns an image sha")
		outputImage = strings.TrimSpace(stdOutBuffer.String())
		Expect(outputImage).To(MatchRegexp("^sha256:[a-f0-9]+$"))

		By("checking if the label exists")
		inspectOutput, _, err := cli.ImageInspectWithRaw(context.TODO(), outputImage)
		Expect(err).ToNot(HaveOccurred())

		labelValue := inspectOutput.Config.Labels["io.pivotal.metadata"]
		Expect(labelValue).To(Equal("metadata here"))

		By("checking that the input image is parent of the output image")
		inspectInput, _, err := cli.ImageInspectWithRaw(context.TODO(), inputImage)
		Expect(err).ToNot(HaveOccurred())

		Expect(inspectOutput.Parent).To(Equal(inspectInput.ID))
	})

	AfterEach(func() {
		_, err := cli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())
	})
})
