package additionalsources

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/image"

	"github.com/pivotal/deplab/pkg/metadata"
)

type HTTPHeadFn func(url string) (resp *http.Response, err error)

func ArchiveUrlProvider(_ image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	for _, archiveURL := range params.AdditionalSourceUrls {
		dependency, err := BuildArchiveDependencyMetadata(archiveURL)
		if err != nil {
			return metadata.Metadata{}, err
		}
		ok, message := IsValidURL(archiveURL, http.Head)
		if !ok {
			errMsg := fmt.Sprintf("failed to validate additional source url: %s", message)
			if params.IgnoreValidationErrors {
				log.Printf("warning: %s", errMsg) //TODO return warning?
			} else {
				return metadata.Metadata{}, fmt.Errorf("error: %s", errMsg)
			}
		}
		md.Dependencies = append(md.Dependencies, dependency)
	}

	return md, nil
}

func AdditionalSourcesProvider(_ image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	var archiveUrls []string
	for _, additionalSourcesFile := range params.AdditionalSourceFilePaths {
		archiveUrlsFromAdditionalSourcesFile, gitVcsFromAdditionalSourcesFile, err := ParseAdditionalSourcesFile(additionalSourcesFile)
		if err != nil {
			errMsg := fmt.Sprintf("could not parse additional sources file: %s, %s", additionalSourcesFile, err)
			if params.IgnoreValidationErrors {
				log.Printf("warning: %s", errMsg) //TODO want to return warnings, and deal with logging in caller
			} else {
				return metadata.Metadata{}, fmt.Errorf("error: %s", errMsg)
			}
		}
		archiveUrls = append(archiveUrls, archiveUrlsFromAdditionalSourcesFile...)
		md.Dependencies = append(md.Dependencies, gitVcsFromAdditionalSourcesFile...)
	}
	return ArchiveUrlProvider(nil, common.RunParams{AdditionalSourceUrls: archiveUrls}, md)
}

func BuildArchiveDependencyMetadata(archiveUrl string) (metadata.Dependency, error) {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: metadata.ArchiveType,
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
	var errorMessages []string
	for _, asu := range additionalSourceUrls {
		if ok, message := IsValidURL(asu, fn); !ok {
			errorMessages = append(errorMessages, message)
		}
	}
	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, ", "))
	}
	return nil
}

func IsValidURL(additionalSourceUrl string, fn HTTPHeadFn) (bool, string) {
	if !isValidExtension(additionalSourceUrl) {
		return false, fmt.Sprintf("unsupported extension for url %s", additionalSourceUrl)
	}

	if resp, err := fn(additionalSourceUrl); err != nil {
		return false, fmt.Sprintf("invalid url: %s", err)
	} else {
		if resp.StatusCode > 299 {
			return false, fmt.Sprintf("got status code %d when trying to reach %s (expected 2xx)", resp.StatusCode, additionalSourceUrl)
		}
	}
	return true, ""
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
