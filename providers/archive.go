package providers

import (
	"net/http"

	"github.com/pivotal/deplab/metadata"
)

type HTTPHeadFn func(url string) (resp *http.Response, err error)

const ArchiveType = "archive"

func BuildArchiveDependencyMetadata(archiveUrl string) (metadata.Dependency, error) {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: ArchiveType,
			Metadata: metadata.ArchiveSourceMetadata{
				URL: archiveUrl,
			},
		},
	}, nil
}

func ValidateURLs(additionalSourceUrls []string, fn HTTPHeadFn) error {
	for _, asu := range additionalSourceUrls {
		if _, err := fn(asu); err != nil {
			return err
		}
	}
	return nil
}
