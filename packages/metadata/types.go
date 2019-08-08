package metadata

type Metadata struct {
	Dependencies []Dependency `json:"dependencies"`
}

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
	Packages []Package `json:"packages"`
}

type Package struct {
	Package      string        `json:"package"`
	Version      string        `json:"version"`
	Architecture string        `json:"architecture"`
	Source       PackageSource `json:"source"`
}

type PackageSource struct {
	Package         string `json:"package"`
	Version         string `json:"version"`
	UpstreamVersion string `json:"upstreamVersion"`
}
