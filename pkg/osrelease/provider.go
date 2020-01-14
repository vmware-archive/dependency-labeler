package osrelease

import (
	"strings"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/pivotal/deplab/pkg/image"

	"github.com/joho/godotenv"
)

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	md.Base = BuildOSMetadata(dli)
	return md, nil
}

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
