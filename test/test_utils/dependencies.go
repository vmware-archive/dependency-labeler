package test_utils

import "github.com/pivotal/deplab/pkg/metadata"

func SelectDpkgDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.DebianPackageListSourceType)
}

func SelectRpmDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.RPMPackageListSourceType)
}

func SelectBuildpackDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	return metadata.SelectDependency(dependencies, metadata.BuildpackMetadataType)
}
