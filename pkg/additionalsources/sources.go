// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package additionalsources

import (
	"fmt"
	"os"
	"strings"

	"github.com/pivotal/deplab/pkg/git"

	"github.com/pivotal/deplab/pkg/metadata"

	"gopkg.in/yaml.v2"
)

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
		case metadata.GitSourceType:
			if !git.IsValidGitDependency(vcs.Url) {
				errorMessages = append(errorMessages, fmt.Sprintf("vcs git url in an unsupported format: %s", vcs.Url))
			}
			gitDependencies = append(gitDependencies, CreateGitDependency(vcs))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("unsupported vcs protocol: %s", vcs.Protocol))
		}
	}

	if len(errorMessages) != 0 {
		return urls, gitDependencies, fmt.Errorf(strings.Join(errorMessages, ", "))
	}

	return urls, gitDependencies, nil
}

func CreateGitDependency(vcs AdditionalSourceVcs) metadata.Dependency {
	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: metadata.GitSourceType,
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
