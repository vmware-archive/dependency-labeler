package providers

import (
	"strings"

	"github.com/pivotal/deplab/rootfs"

	"github.com/joho/godotenv"
	"github.com/pivotal/deplab/metadata"
)

func BuildOSMetadata(dli rootfs.Image) metadata.Base {
	osRelease, err := dli.GetFileContent("/etc/os-release")

	if err != nil {
		return metadata.UnknownBase
	}

	envMap, err := godotenv.Unmarshal(osRelease)
	if err != nil {
		return metadata.UnknownBase
	}

	mdBase := metadata.Base{}
	for k, v := range envMap {
		mdBase[strings.ToLower(k)] = v
	}

	return mdBase
}
