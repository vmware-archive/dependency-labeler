package providers

import (
	"os/exec"

	"github.com/joho/godotenv"

	"github.com/pivotal/deplab/metadata"
)

func BuildOSMetadata(imageName string) metadata.Base {
	osReleaseCmd := exec.Command("docker", "run", "--rm", imageName, "cat", "/etc/os-release")

	out, err := osReleaseCmd.Output()

	if err != nil {
		return metadata.UnknownBase
	}

	envMap, err := godotenv.Unmarshal(string(out))
	if err != nil {
		return metadata.UnknownBase
	}

	return metadata.Base{
		Name:            envMap["NAME"],
		VersionCodename: envMap["VERSION_CODENAME"],
		VersionID:       envMap["VERSION_ID"],
	}
}
