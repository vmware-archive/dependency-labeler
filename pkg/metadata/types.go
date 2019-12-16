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
	Packages   []DpkgPackage `json:"packages"`
	AptSources []string      `json:"apt_sources"`
}

type RpmPackageListSourceMetadata struct {
	Packages []RpmPackage `json:"packages"`
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
