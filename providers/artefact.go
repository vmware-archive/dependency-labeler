package providers

import (
	"gopkg.in/yaml.v2"
	"os"
)

type ArtefactFile struct {
	Blobs []ArtefactFileBlob `yml:blobs`
}

type ArtefactFileBlob struct {
	Url string `yml:url`
}

func ExtractBlobUrlsFromArtefactFile(artefactFilePath string) ([]string, error) {
	artefactFileReader, err := os.Open(artefactFilePath)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(artefactFileReader)
	var artefacts ArtefactFile
	err = decoder.Decode(&artefacts)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, blob := range artefacts.Blobs{
		urls = append(urls, blob.Url)
	}

	return urls, nil
}
