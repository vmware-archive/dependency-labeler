package providers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"

	"github.com/pivotal/deplab/metadata"
)

func BuildOSMetadata(imageName string) metadata.Base {
	createContainerCmd := exec.Command("docker", "create", imageName, "foo")
	containerId, err := createContainerCmd.Output()
	if err != nil {
		log.Fatalf("failed to create container: %s", err)
	}

	osReleaseCmd1 := exec.Command("docker", "cp", "-L", fmt.Sprintf("%s:%s", strings.TrimSpace(string(containerId)), "/etc/os-release"), "-")
	osReleaseCmd2 := exec.Command("tar", "-xO")

	r, w := io.Pipe()
	osReleaseCmd1.Stdout = w
	osReleaseCmd2.Stdin = r

	var osReleaseOutput2 bytes.Buffer
	osReleaseCmd2.Stdout = &osReleaseOutput2

	err = osReleaseCmd1.Start()
	if err != nil {
		return metadata.UnknownBase
	}

	err = osReleaseCmd2.Start()
	if err != nil {
		return metadata.UnknownBase
	}

	err = osReleaseCmd1.Wait()
	if err != nil {
		return metadata.UnknownBase
	}

	err = w.Close()
	if err != nil {
		return metadata.UnknownBase
	}

	err = osReleaseCmd2.Wait()
	if err != nil {
		return metadata.UnknownBase
	}

	envMap, err := godotenv.Unmarshal(osReleaseOutput2.String())
	if err != nil {
		return metadata.UnknownBase
	}

	if envMap["NAME"] == "" {
		return metadata.UnknownBase
	}

	mdbase := metadata.Base{}
	for k, v := range envMap {
		mdbase[strings.ToLower(k)] = v
	}

	return mdbase
}
