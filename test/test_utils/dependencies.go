// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package test_utils

import "github.com/vmware-tanzu/dependency-labeler/pkg/metadata"

func SelectDpkgDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.DebianPackageListSourceType)
}

func SelectRpmDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.RPMPackageListSourceType)
}

func SelectBuildpackDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.BuildpackMetadataType)
}

func SelectKpackDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.PackageType)
}
