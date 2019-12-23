package common

import (
	"github.com/pivotal/deplab/pkg/image"
	"github.com/pivotal/deplab/pkg/metadata"
)

type Provider func(image image.Image) (metadata.Dependency, error)
