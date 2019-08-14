package providers

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"

	"github.com/pivotal/deplab/metadata"
)

func BuildOSMetadata(imageName string) (*metadata.Base, error) {
	osReleaseCmd := exec.Command("docker", "run", "--rm", imageName, "cat", "/etc/os-release")

	out, err := osReleaseCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cannot get os-release: %s", err)
	}

	envMap, err := godotenv.Unmarshal(string(out))
	if err != nil {
		return nil, fmt.Errorf("cannot parse os-release: %s", err)
	}

	return &metadata.Base{
		Name:            envMap["NAME"],
		VersionCodename: envMap["VERSION_CODENAME"],
		VersionID:       envMap["VERSION_ID"],
	}, nil
}
