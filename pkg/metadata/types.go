// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package metadata

const (
	DebianPackageListSourceType = "debian_package_list"
	GitSourceType               = "git"
	RPMPackageListSourceType    = "rpm_package_list"
	ArchiveType                 = "archive"
	PackageType                 = "package"
	BuildpackMetadataType       = "buildpack_metadata"
)

type Metadata struct {
	Base         Base         `json:"base"`
	Provenance   []Provenance `json:"provenance"`
	Dependencies []Dependency `json:"dependencies"`
}

type Provenance struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

type Base map[string]string

type Dependency struct {
	Type   string `json:"type"`
	Source Source `json:"source"`
}

type Source struct {
	Type     string                 `json:"type"`
	Version  map[string]interface{} `json:"version"`
	Metadata interface{}            `json:"metadata"`
}

type DebianPackageListSourceMetadata struct {
	Packages   []DpkgPackage `json:"packages"`
	AptSources []string      `json:"apt_sources"`
}

type RpmPackageListSourceMetadata struct {
	Packages []RpmPackage `json:"packages"`
}

type BuildpackBOMSourceMetadata struct {
	Buildpacks      []Buildpack            `json:"buildpacks"`
	BillOfMaterials []BuildpackBOM         `json:"bom"`
	Launcher        map[string]interface{} `json:"launcher"`
}

type KpackRepoSourceMetadata struct {
	Url  string        `json:"url"`
	Refs []interface{} `json:"refs"`
}

type GitSourceMetadata struct {
	URL  string   `json:"url"`
	Refs []string `json:"refs"`
}

type ArchiveSourceMetadata struct {
	URL string `json:"url"`
}

type DpkgPackage struct {
	Package      string        `json:"package"`
	Version      string        `json:"version"`
	Architecture string        `json:"architecture"`
	Source       PackageSource `json:"source"`
}

type RpmPackage struct {
	Package      string `json:"package" rpm:"NAME"`
	Version      string `json:"version" rpm:"VERSION"`
	Architecture string `json:"architecture" rpm:"ARCH"`
	License      string `json:"license" rpm:"LICENSE"`
	SourceRpm    string `json:"source_rpm" rpm:"SOURCERPM"`
}

type Buildpack struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type BuildpackBOM struct {
	Name      string               `json:"name"`
	Version   string               `json:"version"`
	Metadata  BuildpackBOMMetadata `json:"metadata"`
	Buildpack Buildpack            `json:"buildpack"`
}

type PackageSource struct {
	Package         string `json:"package"`
	Version         string `json:"version"`
	UpstreamVersion string `json:"upstreamVersion"`
}

type BuildpackBOMMetadata map[string]interface{}

var UnknownBase = Base{
	"name":             "unknown",
	"version_codename": "unknown",
	"version_id":       "unknown",
}

var BusyboxBase = Base{
	"name":             "busybox",
	"pretty_name":      "busybox",
	"version_codename": "unknown",
	"version_id":       "unknown",
}

var ScratchBase = Base{
	"name":             "scratch",
	"pretty_name":      "scratch",
	"version_codename": "unknown",
	"version_id":       "unknown",
}
