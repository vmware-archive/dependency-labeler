// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package osrelease

import (
	"strings"

	"github.com/vmware-tanzu/dependency-labeler/pkg/common"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"

	"github.com/vmware-tanzu/dependency-labeler/pkg/image"

	"github.com/joho/godotenv"
)

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	md.Base = BuildOSMetadata(dli)
	return md, nil
}

func BuildOSMetadata(dli image.Image) metadata.Base {
	osRelease, err := dli.GetFileContent("/etc/os-release")

	if err != nil {
		binContents, binErr := dli.GetDirFileNames("/bin", false)
		if binErr != nil {
			return checkScratchBase(dli)
		}

		hasAsh := contains(binContents, "ash")
		if hasAsh {
			return metadata.BusyboxBase
		}

		return checkScratchBase(dli)
	}

	envMap, err := godotenv.Unmarshal(osRelease)
	if err != nil {
		return metadata.UnknownBase
	}

	mdBase := metadata.Base{}
	for k, v := range envMap {
		mdBase[strings.ToLower(k)] = v
	}

	return mdBase
}

func checkScratchBase(dli image.Image) metadata.Base {
	rootContents, rootErr := dli.GetDirFileNames("/", true)
	if rootErr == nil  && len(rootContents) <= 3 {
		return metadata.ScratchBase
	}

	return metadata.UnknownBase
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
