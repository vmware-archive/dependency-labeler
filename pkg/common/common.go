package common

import (
	"github.com/pivotal/deplab/pkg/image"
	"github.com/pivotal/deplab/pkg/metadata"
)

type Provider interface {
	BuildDependencyMetadata(image image.Image) (metadata.Dependency, error)
}
