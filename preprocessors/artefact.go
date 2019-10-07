package preprocessors

import (
	"github.com/pivotal/deplab/metadata"
	"gopkg.in/yaml.v2"
	"os"
)

type ArtefactFile struct {
	Blobs []ArtefactFileBlob `yml:blobs`
	Vcs []ArtefactFileVcs `yml:vcs`
}

type ArtefactFileBlob struct {
	Url string `yml:url`
}

type ArtefactFileVcs struct {
	Protocol string `yml:protocol`
	Version string `yml:version`
	Url string `yml:url`
	Refs string `yml:refs`
}

func ParseArtefactFile(artefactFilePath string) ([]string, []metadata.Dependency, error) {
	artefactFileReader, err := os.Open(artefactFilePath)
	if err != nil {
		return nil, nil, err
	}

	decoder := yaml.NewDecoder(artefactFileReader)
	var artefacts ArtefactFile
	err = decoder.Decode(&artefacts)
	if err != nil {
		return nil, nil, err
	}

	var urls []string
	for _, blob := range artefacts.Blobs{
		urls = append(urls, blob.Url)
	}

	var vcsArtefacts []metadata.Dependency
	for _, vcs := range artefacts.Vcs{
		if vcs.Protocol == "git"{
			vcsArtefacts = append(vcsArtefacts, createGitDependency(vcs))
		}
	}

	return urls, vcsArtefacts, nil
}

func createGitDependency(vcs ArtefactFileVcs) metadata.Dependency {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: "git",
			Version: map[string]interface{}{
				"commit": vcs.Version,
			},
			Metadata: metadata.GitSourceMetadata{
				URL:  vcs.Url,
				Refs: []string{vcs.Refs},
			},
		},
	}
}
