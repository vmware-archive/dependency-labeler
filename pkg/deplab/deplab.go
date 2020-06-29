// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package deplab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pivotal/deplab/pkg/kpack"
	"os"
	"strings"

	"github.com/pivotal/deplab/pkg/cnb"

	"github.com/pivotal/deplab/pkg/additionalsources"
	"github.com/pivotal/deplab/pkg/common"
	"github.com/pivotal/deplab/pkg/rpm"

	"github.com/pivotal/deplab/pkg/git"

	"github.com/pivotal/deplab/pkg/dpkg"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/pivotal/deplab/pkg/osrelease"

	"github.com/pivotal/deplab/pkg/image"
)

type provider func(image.Image, common.RunParams, metadata.Metadata) (metadata.Metadata, error)

var Version = "0.0.0-dev"
var Provenance = metadata.Provenance{
	Name:    "deplab",
	Version: Version,
	URL:     "https://github.com/pivotal/deplab",
}

func Run(params common.RunParams) error {
	dli, err := image.NewDeplabImage(params.InputImage, params.InputImageTarPath)

	if err != nil {
		return fmt.Errorf("could not load image: %w", err)
	}
	defer dli.Cleanup()

	md := metadata.Metadata{Dependencies: make([]metadata.Dependency, 0)}

	for _, provider := range []provider{
		dpkg.Provider,
		rpm.Provider,
		cnb.Provider,
		git.Provider,
		additionalsources.ArchiveUrlProvider,
		additionalsources.AdditionalSourcesProvider,
		osrelease.Provider,
		ProvenanceProvider,
	} {
		if md2, err := provider(&dli, params, md); err == nil {
			md = md2
		} else {
			return fmt.Errorf("error generating dependencies: %w", err)
		}
	}

	err = writeOutputs(dli, params, md)
	if err != nil {
		return fmt.Errorf("could not write outputs: %w", err)
	}

	return nil
}

func RunInspect(inputImage, inputImageTar string) error {
	dli, err := image.NewDeplabImage(inputImage, inputImageTar)

	if err != nil {
		return fmt.Errorf("inspect cannot open the provided image from '%s%s': %s", inputImage, inputImageTar, err)
	}

	inspectMetadata := metadata.Metadata{}

	for _, provider := range []provider{
		dpkg.Provider,
		rpm.Provider,
		cnb.Provider,
		osrelease.Provider,
		kpack.Provider,
		ProvenanceProvider,
		ExistingLabelProvider,
	} {
		if md2, err := provider(&dli, common.RunParams{}, inspectMetadata); err == nil {
			inspectMetadata = md2
		} else {
			return fmt.Errorf("inspect error generating dependencies for image '%s%s': %w", inputImageTar, inputImage, err)
		}
	}

	label, err := json.Marshal(inspectMetadata)
	if err != nil {
		return fmt.Errorf("cannot generate json: %w", err)
	}

	stdOutBuffer := bytes.Buffer{}
	err = json.Indent(&stdOutBuffer, label, "", "  ")
	if err != nil {
		return fmt.Errorf("inspect cannot pretty print the label of the provided image '%s%s', label: %s: %w", inputImageTar, inputImage, label, err)
	}

	fmt.Println(stdOutBuffer.String())
	return nil
}

func writeOutputs(dli image.Image, params common.RunParams, md metadata.Metadata) error {
	if params.OutputImageTar != "" {
		err := dli.ExportWithMetadata(md, params.OutputImageTar, params.Tag)

		if err != nil {
			return fmt.Errorf("error exporting tar to %s: %w", params.OutputImageTar, err)
		}
	}

	if params.MetadataFilePath != "" {
		err := metadata.WriteMetadataFile(md, params.MetadataFilePath)
		if err != nil {
			return fmt.Errorf("could not write metadata file: %w", err)
		}
	}

	if params.DpkgFilePath != "" {
		err := dpkg.WriteDpkgFile(md, params.DpkgFilePath, Version)
		if err != nil {
			return fmt.Errorf("could not write dpkg file: %w", err)
		}
	}

	return nil
}

func ProvenanceProvider(_ image.Image, _ common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	md.Provenance = append(md.Provenance, Provenance)
	return md, nil
}

func ExistingLabelProvider(dli image.Image, _ common.RunParams, md metadata.Metadata) (m metadata.Metadata, err error) {
	cf, err := dli.GetConfig()
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("cannot retrieve the Config file: %w", err)
	}

	existingMetadata := metadata.Metadata{}
	existinglabel, foundDeplab := cf.Config.Labels["io.deplab.metadata"]

	if foundDeplab {
		var err = json.Unmarshal([]byte(existinglabel), &existingMetadata)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("cannot parse the label %s: %w", existinglabel, err)
		}
	} else {
		if existinglabel, ok := cf.Config.Labels["io.pivotal.metadata"]; ok {
			var err = json.Unmarshal([]byte(existinglabel), &existingMetadata)
			if err != nil {
				return metadata.Metadata{}, fmt.Errorf("cannot parse the label %s: %w", existinglabel, err)
			}
		}
	}

	mergedMetadata, warnings := metadata.Merge(existingMetadata, md)
	if len(warnings) > 0 {
		var warnStrings []string
		for _, warning := range warnings {
			warnStrings = append(warnStrings, string(warning))
		}
		fmt.Fprintln(os.Stderr, "Difference identified in matching metadata elements already present on image:", strings.Join(warnStrings, ", "))
	}

	return mergedMetadata, nil
}
