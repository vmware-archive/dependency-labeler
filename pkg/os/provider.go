package os

import (
	"strings"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/pivotal/deplab/pkg/image"

	"github.com/joho/godotenv"
)

func BuildOSMetadata(dli image.Image) metadata.Base {
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
