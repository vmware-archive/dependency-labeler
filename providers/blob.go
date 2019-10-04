package providers

import (
	"github.com/pivotal/deplab/metadata"
)

func BuildBlobDependencyMetadata(blobUrl string) (metadata.Dependency, error) {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: "blob",
			Metadata: metadata.BlobSourceMetadata{
				URL:  blobUrl,
			},
		},
	}, nil
}
