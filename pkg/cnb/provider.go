// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package cnb

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/vmware-tanzu/dependency-labeler/pkg/common"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/vmware-tanzu/dependency-labeler/pkg/image"
)

func Provider(dli image.Image, _ common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	var buildpackMetadataContents string
	config, err := dli.GetConfig()

	if err != nil {
		return metadata.Metadata{}, err
	}

	buildpackMetadataContents = config.Config.Labels["io.buildpacks.build.metadata"]

	if buildpackMetadataContents != "" {
		buildpackMetadata, err := parseMetadataJSON(buildpackMetadataContents)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("could not parse buildpack metadata toml: %w", err)
		}

		version, err := common.Digest(buildpackMetadata)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("could not get digest for buildpack metadata: %w", err)
		}

		md.Dependencies = append(md.Dependencies, metadata.Dependency{
			Type: metadata.BuildpackMetadataType,
			Source: metadata.Source{
				Type:     "inline",
				Metadata: buildpackMetadata,
				Version: map[string]interface{}{
					"sha256": version,
				},
			},
		})
	}
	return md, nil
}

func parseMetadataJSON(buildpackMetadata string) (metadata.BuildpackBOMSourceMetadata, error) {
	var bp metadata.BuildpackBOMSourceMetadata

	err := json.Unmarshal([]byte(buildpackMetadata), &bp)
	if err != nil {
		return metadata.BuildpackBOMSourceMetadata{}, fmt.Errorf("could not decode json: %w", err)
	}

	collator := collate.New(language.BritishEnglish)
	sort.Slice(bp.Buildpacks, func(i, j int) bool {
		return collator.CompareString(bp.Buildpacks[i].ID, bp.Buildpacks[j].ID) < 0
	})

	sort.Slice(bp.BillOfMaterials, func(i, j int) bool {
		return collator.CompareString(bp.BillOfMaterials[i].Name, bp.BillOfMaterials[j].Name) < 0
	})

	return bp, nil
}
