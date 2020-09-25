// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package kpack

import (
	"encoding/json"
	"fmt"
	"github.com/vmware-tanzu/dependency-labeler/pkg/common"
	"github.com/vmware-tanzu/dependency-labeler/pkg/image"
	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

type RepoSource struct {
	Source Source `json:"source"`
}

type Source struct {
	Type     string            `json:"type"`
	Version  map[string]string `json:"version"`
	Metadata map[string]string `json:"metadata"`
}

func Provider(dli image.Image, _ common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	var kpackMetadataContents string

	config, err := dli.GetConfig()
	if err != nil {
		return metadata.Metadata{}, err
	}

	kpackMetadataContents = config.Config.Labels["io.buildpacks.project.metadata"]

	if kpackMetadataContents != "" {
		dep, err := parseMetadataJSON(kpackMetadataContents)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("could not parse kpack metadata: %w", err)
		}

		md.Dependencies = append(md.Dependencies, dep)
	}

	return md, nil
}

func parseMetadataJSON(kpackMetadata string) (metadata.Dependency, error) {
	var md RepoSource

	err := json.Unmarshal([]byte(kpackMetadata), &md)
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("could not decode json: %w", err)
	}

	var kpackMd = metadata.KpackRepoSourceMetadata{
		Url:  md.Source.Metadata["repository"],
		Refs: []interface{}{},
	}

	var dep = metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: "git",
			Version: map[string]interface{}{
				"commit": md.Source.Version["commit"],
			},
			Metadata: kpackMd,
		},
	}
	return dep, nil
}
