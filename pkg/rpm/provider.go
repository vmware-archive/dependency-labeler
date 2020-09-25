// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package rpm

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/vmware-tanzu/dependency-labeler/pkg/common"

	"github.com/vmware-tanzu/dependency-labeler/pkg/image"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

const RPMDbPath = "/var/lib/rpm"

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {

	absPath, err := dli.AbsolutePath(RPMDbPath)
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("absolute path for rpm database: %w", err)
	}

	exists, err := exists(path.Join(absPath, "Packages"))
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("rpm could not find existance of path: %w", err)
	}
	if !exists {
		return md, nil
	}

	if !isRPMInstalled() {
		return metadata.Metadata{}, fmt.Errorf("an rpm database exists at %s but rpm is not installed and available on your path: %w", RPMDbPath, err)
	}

	query := QueryFormat()
	cmd := exec.Command("rpm",
		"-qa",
		"--dbpath", absPath,
		"--queryformat", query,
	)
	stdOutBuffer := &strings.Builder{}
	cmd.Stdout = stdOutBuffer

	err = cmd.Run()

	if err != nil {
		return metadata.Metadata{},
			fmt.Errorf("failed to execute rpm at path, %s, with query, %s: %w", absPath, query, err)
	}

	if strings.TrimSpace(stdOutBuffer.String()) == "" {
		return metadata.Metadata{}, fmt.Errorf("no rpm packages data found")
	}

	allPackagesDetails := strings.Split(strings.TrimSpace(stdOutBuffer.String()), "\n")

	var packages []metadata.RpmPackage

	for _, line := range allPackagesDetails {
		packages = append(packages, UnmarshalPackage(line))
	}
	collator := collate.New(language.BritishEnglish)
	sort.Slice(packages, func(i, j int) bool {
		return collator.CompareString(packages[i].Package, packages[j].Package) < 0
	})

	sourceMetadata := metadata.RpmPackageListSourceMetadata{
		Packages: packages,
	}

	version, err := common.Digest(sourceMetadata)
	if err != nil {
		return metadata.Metadata{}, fmt.Errorf("could not get digest for source metadata: %w", err)
	}

	md.Dependencies = append(md.Dependencies, metadata.Dependency{
		Type: metadata.RPMPackageListSourceType,
		Source: metadata.Source{
			Type: "inline",
			Version: map[string]interface{}{
				"sha256": version,
			},
			Metadata: sourceMetadata,
		},
	})

	return md, nil
}

func isRPMInstalled() bool {
	stdOutBuffer := &strings.Builder{}
	cmd := exec.Command("rpm",
		"--version",
	)

	cmd.Stdout = stdOutBuffer

	err := cmd.Run()

	return err == nil && strings.Contains(stdOutBuffer.String(), "RPM version")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
