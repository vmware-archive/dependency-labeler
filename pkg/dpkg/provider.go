// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package dpkg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/pivotal/deplab/pkg/image"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

//func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
//	dependency, err := BuildDependencyMetadata(dli)
//	if err != nil {
//		return metadata.Metadata{}, err
//	}
//	md.Dependencies = append(md.Dependencies, dependency)
//	return md, nil
//}

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	packages := getDebianPackages(dli)

	if len(packages) != 0 {
		sources, err := getAptSources(dli)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("could not get apt sources: %w", err)
		}

		sourceMetadata := metadata.DebianPackageListSourceMetadata{
			Packages:   packages,
			AptSources: sources,
		}

		version, err := common.Digest(sourceMetadata)
		if err != nil {
			return metadata.Metadata{}, fmt.Errorf("could not get digest for source metadata: %w", err)
		}

		md.Dependencies = append(md.Dependencies, metadata.Dependency{
			Type: metadata.DebianPackageListSourceType,
			Source: metadata.Source{
				Type: "inline",
				Version: map[string]interface{}{
					"sha256": version,
				},
				Metadata: sourceMetadata,
			},
		})

	}
	return md, nil
}

func getAptSources(dli image.Image) ([]string, error) {
	sources := []string{}
	fileListContent, err := dli.GetDirContents("/etc/apt/sources.list.d")
	if err != nil {
		// in this case an empty or non-existant directory is not an error
		fileListContent = []string{}
	}

	aFileListContent, err := dli.GetFileContent("/etc/apt/sources.list")
	if err == nil {
		// in this case an empty or non-existant file is not an error
		fileListContent = append(fileListContent, aFileListContent)
	}

	for _, fileContent := range fileListContent {
		for _, content := range strings.Split(fileContent, "\n") {
			trimmed := strings.TrimSpace(content)
			if !strings.HasPrefix(trimmed, "#") && len(trimmed) != 0 {
				sources = append(sources, content)
			}
		}
	}

	collator := collate.New(language.BritishEnglish)
	sort.Slice(sources, func(i, j int) bool {
		return collator.CompareString(sources[i], sources[j]) < 0
	})

	return sources, nil
}

func getDebianPackages(dli image.Image) ([]metadata.DpkgPackage) {
	var packages []metadata.DpkgPackage

	packages = append(packages, listPackagesFromStatus(dli)...)
	packages = append(packages, listPackagesFromStatusD(dli)...)

	collator := collate.New(language.BritishEnglish)
	sort.Slice(packages, func(i, j int) bool {
		return collator.CompareString(packages[i].Package, packages[j].Package) < 0
	})

	return packages
}

func ParseStatDBEntry(content string) (metadata.DpkgPackage, error) {
	pkg := metadata.DpkgPackage{}

	if strings.TrimSpace(content) == "" {
		return pkg, fmt.Errorf("invalid StatDB entry")
	}

	for _, inputLine := range strings.Split(content, "\n") {
		idx := strings.Index(inputLine, ":")
		if idx == -1 {
			continue
		}
		key := inputLine[0:idx]
		value := strings.TrimSpace(inputLine[idx+1:])
		switch key {
		case "Package":
			pkg.Package = value
		case "Version":
			pkg.Version = value
		case "Architecture":
			pkg.Architecture = value
		case "Source":
			idx := strings.Index(value, "(")
			if idx == -1 {
				pkg.Source.Package = strings.TrimSpace(value)
			} else {
				pkg.Source.Package = strings.TrimSpace(value[0:idx])
				version := strings.Trim(value[idx:], " ()")
				pkg.Source.Version = version
				pkg.Source.UpstreamVersion = getUpstreamVersion(version)
			}
		default:
			continue
		}
	}

	if pkg.Source.Package == "" {
		pkg.Source = metadata.PackageSource{
			Package:         pkg.Package,
			Version:         pkg.Version,
			UpstreamVersion: getUpstreamVersion(pkg.Version),
		}
	}
	if pkg.Source.Version == "" {
		pkg.Source.Version = pkg.Version
		pkg.Source.UpstreamVersion = getUpstreamVersion(pkg.Version)
	}

	return pkg, nil
}

func listPackagesFromStatusD(dli image.Image) (packages []metadata.DpkgPackage) {
	fileList, err := dli.GetDirContents("/var/lib/dpkg/status.d")
	if err != nil {
		// in this case an empty or non-existent directory is not an error
		fileList = []string{}
	}

	for _, file := range fileList {
		packageEntry, err := ParseStatDBEntry(file)
		if err == nil {
			packages = append(packages, packageEntry)
		}
	}

	return packages
}

func listPackagesFromStatus(dli image.Image) (packages []metadata.DpkgPackage) {
	statDBString, err := dli.GetFileContent("/var/lib/dpkg/status")
	if err != nil {
		// in this case an empty or non-existent file is not an error
		statDBString = ""
	}

	statDBEntries := strings.Split(statDBString, "\n\n")
	for _, entryString := range statDBEntries {
		entry, err := ParseStatDBEntry(entryString)
		if err == nil {
			packages = append(packages, entry)
		}
	}

	return packages
}

func getUpstreamVersion(input string) string {
	version := strings.Split(input, "-")[0]
	if strings.Contains(version, ":") {
		version = strings.Split(version, ":")[1]
	}
	return version
}
