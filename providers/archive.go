package providers

import (
	"fmt"
	"net/http"
	"strings"

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

// https://en.wikipedia.org/wiki/Tar_(computing)#Suffixes_for_compressed_files
var SupportedExtensions = []string{
	"7z",
	"tar.bz2",
	"tar.gz",
	"tar.lz",
	"tar.lzma",
	"tar.lzo",
	"tar.xz",
	"tar.Z",
	"tar.zst",
	"taz",
	"taZ",
	"tb2",
	"tbz",
	"tbz2",
	"tgz",
	"tlz",
	"tpz",
	"txz",
	"tZ",
	"tz2",
	"tzst",
	"tar.bz2",
	"zip",
}

func ValidateURLs(additionalSourceUrls []string, fn HTTPHeadFn) error {
	for _, asu := range additionalSourceUrls {
		if !isValidExtension(asu) {
			return fmt.Errorf("unsupported extension for url %s", asu)
		}

		if resp, err := fn(asu); err != nil {
			return fmt.Errorf("invalid url: %s", err)
		} else {
			if resp.StatusCode > 299 {
				return fmt.Errorf("got status code %d when trying to reach %s (expected 2xx)", resp.StatusCode, asu)
			}
		}
	}
	return nil
}

func isValidExtension(sourceUrl string) bool {
	for _, extension := range SupportedExtensions {
		if strings.HasSuffix(sourceUrl, "."+extension) ||
			strings.Contains(sourceUrl, "."+extension+"#") {
			return true
		}
	}
	return false
}
