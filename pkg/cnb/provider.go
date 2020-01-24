package cnb

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/metadata"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/pivotal/deplab/pkg/image"
)

func Provider(dli image.Image, _ common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	dependency, err := buildDependencyMetadata(dli)
	if err != nil {
		return metadata.Metadata{}, err
	}
	md.Dependencies = append(md.Dependencies, dependency)
	return md, nil
}

func buildDependencyMetadata(dli image.Image) (metadata.Dependency, error) {
	var buildpackMetadataContents string
	config, err := dli.GetConfig()

	if err != nil {
		return metadata.Dependency{}, err
	}

	buildpackMetadataContents = config.Config.Labels["io.buildpacks.build.metadata"]

	if buildpackMetadataContents != "" {
		buildpackMetadata, err := parseMetadataJSON(buildpackMetadataContents)
		if err != nil {
			return metadata.Dependency{}, fmt.Errorf("could not parse buildpack metadata toml: %w", err)
		}

		version, err := common.Digest(buildpackMetadata)
		if err != nil {
			return metadata.Dependency{}, fmt.Errorf("could not get digest for buildpack metadata: %w", err)
		}

		return metadata.Dependency{
			Type: metadata.BuildpackMetadataType,
			Source: metadata.Source{
				Type:     "inline",
				Metadata: buildpackMetadata,
				Version: map[string]interface{}{
					"sha256": version,
				},
			},
		}, nil
	}

	return metadata.Dependency{}, nil
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
