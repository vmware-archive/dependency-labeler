package metadata

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
	Packages   []Package `json:"packages"`
	AptSources []string  `json:"apt_sources"`
}

type GitSourceMetadata struct {
	URL  string   `json:"url"`
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

var UnknownBase = Base{
	"name":             "unknown",
	"version_codename": "unknown",
	"version_id":       "unknown",
}
