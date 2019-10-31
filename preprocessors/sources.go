package preprocessors

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pivotal/deplab/metadata"
	"gopkg.in/yaml.v2"
)

type AdditionalSources struct {
	Archives []AdditionalSourceArchive `yml:archives`
	Vcs      []AdditionalSourceVcs     `yml:vcs`
}

type AdditionalSourceArchive struct {
	Url string `yml:url`
}

type AdditionalSourceVcs struct {
	Protocol string `yml:protocol`
	Version  string `yml:version`
	Url      string `yml:url`
}

func ParseAdditionalSourcesFile(additionalSourcesFilePath string) ([]string, []metadata.Dependency, error) {
	additionalSourcesFileReader, err := os.Open(additionalSourcesFilePath)
	if err != nil {
		return nil, nil, err
	}

	decoder := yaml.NewDecoder(additionalSourcesFileReader)
	var additionalSources AdditionalSources
	err = decoder.Decode(&additionalSources)
	if err != nil {
		return nil, nil, err
	}

	var urls []string
	for _, archive := range additionalSources.Archives {
		urls = append(urls, archive.Url)
	}

	var gitDependencies []metadata.Dependency
	var errorMessages []string
	for _, vcs := range additionalSources.Vcs {
		switch vcs.Protocol {
		case GitSourceType:
			if !validGitDependency(vcs.Url) {
				errorMessages = append(errorMessages, fmt.Sprintf("vcs git url in an unsupported format: %s", vcs.Url))
			}
			gitDependencies = append(gitDependencies, createGitDependency(vcs))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("unsupported vcs protocol: %s", vcs.Protocol))
		}
	}

	if len(errorMessages) != 0 {
		return urls, gitDependencies, errors.New(strings.Join(errorMessages, ", "))
	}

	return urls, gitDependencies, nil
}

func validGitDependency(gitUrl string) bool {
	valid, err := regexp.MatchString(`((git|ssh|http(s)?)|(git@[\w\.]+))(:)([\w\.@\:/\-~]+)(/)?`, gitUrl)
	if err != nil {
		log.Printf("error when matching regex to validate git dependency: %s", err)
	}
	return valid
}

func createGitDependency(vcs AdditionalSourceVcs) metadata.Dependency {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: GitSourceType,
			Version: map[string]interface{}{
				"commit": vcs.Version,
			},
			Metadata: metadata.GitSourceMetadata{
				URL:  vcs.Url,
				Refs: []string{},
			},
		},
	}
}
