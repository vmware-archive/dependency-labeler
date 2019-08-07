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

func runDepLab(inputImage []string, expErrCode int) (stdOutBuffer bytes.Buffer, stdErrBuffer bytes.Buffer) {
	stdOutBuffer = bytes.Buffer{}
	stdErrBuffer = bytes.Buffer{}

	cmd := exec.Command(pathToBin, inputImage...)

	session, err := gexec.Start(cmd, &stdOutBuffer, &stdErrBuffer)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session, time.Second*5).Should(gexec.Exit(expErrCode))
	<-session.Exited

	return stdOutBuffer, stdErrBuffer
}

var _ = Describe("deplab", func() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"), client.FromEnv)
	if err != nil {
		panic(err)
	}

	It("labels an image and returns the new sha", func() {

		By("executing it")
		inputImage := "alpine"
		stdOutBuffer, _ := runDepLab([]string{"--image", inputImage}, 0)

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

		_, err = cli.ImageRemove(context.TODO(), outputImage, types.ImageRemoveOptions{})
		Expect(err).ToNot(HaveOccurred())

	})

	It("throws an error if an invalid image sent to docker engine", func() {

		By("executing it")
		inputImage := "swkichtlsmhasd" // random string unlikely for an image ever to exist
		_, stdErrBuffer := runDepLab([]string{"--image", inputImage}, 1)

		errorOutput := strings.TrimSpace(stdErrBuffer.String())
		Expect(errorOutput).To(ContainSubstring("error building image: pull access denied for swkichtlsmhasd, repository does not exist or may require 'docker login'"))
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
})
