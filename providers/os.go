package providers

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/pivotal/deplab/docker"

	"github.com/pivotal/deplab/metadata"
)

func BuildOSMetadata(imageName string) metadata.Base {
	t, err := docker.ReadFromImage(imageName, "/etc/os-release")

	if err != nil {
		return metadata.UnknownBase
	}

	_, err = t.Next()
	if err != nil {
		return metadata.UnknownBase
	}

	envMap, err := godotenv.Parse(t)
	if err != nil {
		return metadata.UnknownBase
	}

	mdBase := metadata.Base{}
	for k, v := range envMap {
		mdBase[strings.ToLower(k)] = v
	}

	return mdBase
}
