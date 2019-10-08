package preprocessors

import (
	"github.com/pivotal/deplab/metadata"
	"gopkg.in/yaml.v2"
	"os"
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
	Version string `yml:version`
	Url string `yml:url`
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
	for _, vcs := range additionalSources.Vcs{
		if vcs.Protocol == "git"{
			gitDependencies = append(gitDependencies, createGitDependency(vcs))
		}
	}

	return urls, gitDependencies, nil
}

func createGitDependency(vcs AdditionalSourceVcs) metadata.Dependency {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: "git",
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
