package metadata

type Metadata struct {
	Base         *Base        `json:"base"`
	Dependencies []Dependency `json:"dependencies"`
}

type Base struct {
	Name            string `json:"name"`
	VersionID       string `json:"version_id"`
	VersionCodename string `json:"version_codename"`
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
	Packages []Package `json:"pkg"`
}

type GitSourceMetadata struct {
	URI  string   `json:"uri"`
	Refs []string `json:"refs"`
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
