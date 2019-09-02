package providers

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/pivotal/deplab/metadata"
)

func BuildDebianDependencyMetadata(imageName string) (metadata.Dependency, error) {
	packages, err := getDebianPackages(imageName)

	if len(packages) != 0 {
		sources, _ := getAptSources(imageName)

		sourceMetadata := metadata.DebianPackageListSourceMetadata{
			Packages:   packages,
			AptSources: sources,
		}

		version := Digest(sourceMetadata)

		dpkgList := metadata.Dependency{
			Type: "debian_package_list",
			Source: metadata.Source{
				Type: "inline",
				Version: map[string]interface{}{
					"sha256": version,
				},
				Metadata: sourceMetadata,
			},
		}

		return dpkgList, nil
	}

	return metadata.Dependency{}, err
}

func Digest(sourceMetadata metadata.DebianPackageListSourceMetadata) string {
	hash := sha256.New()
	encoder := json.NewEncoder(hash)
	_ = encoder.Encode(sourceMetadata)
	version := hex.EncodeToString(hash.Sum(nil))
	return version
}

func getAptSources(imageName string) ([]string, error) {
	stdout := &bytes.Buffer{}

	grep := exec.Command("docker", "run", "--rm", imageName,
		"grep",
		"^[^#]",
		"/etc/apt/sources.list",
		"/etc/apt/sources.list.d",
		"--no-filename",
		"--no-message",
		"--recursive")

	grep.Stdout = stdout

	_ = grep.Run()

	//this requires an empty slice not a nil slice due to JSON serialization
	//nil slices serialize as null
	//empty slice serialize to []
	sources := []string{}

	for _, source := range strings.Split(stdout.String(), "\n") {
		if strings.TrimSpace(source) != "" {
			sources = append(sources, source)
		}
	}

	return sources, nil
}

func getStatus(containerId string) (io.Reader, error) {
	t, err := getTar(containerId, "/var/lib/dpkg/status")
	if err != nil {
		return nil, fmt.Errorf("error retrieving tar from: %s", err)
	}

	if t == nil {
		return nil, nil
	}

	_, err = t.Next()
	if err != nil {
		return nil, fmt.Errorf("error reading status: %s", err)
	}

	return t, nil
}

func getStatusD(containerId string) (*tar.Reader, error) {
	t, err := getTar(containerId, "/var/lib/dpkg/status.d")
	if err != nil {
		return nil, fmt.Errorf("error retrieving tar from: %s", err)
	}

	return t, nil
}

func getTar(containerId string, path string) (*tar.Reader, error) {
	c := exec.Command("docker", "cp", "-L", fmt.Sprintf("%s:%s", strings.TrimSpace(string(containerId)), path), "-")
	var cOut bytes.Buffer
	c.Stdout = &cOut

	err := c.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting docker cp command: %s", err)
	}

	err = c.Wait()
	if err != nil {
		return nil, nil
	}

	return tar.NewReader(&cOut), nil
}

func getDebianPackages(imageName string) ([]metadata.Package, error) {
	createContainerCmd := exec.Command("docker", "create", imageName, "foo")
	containerId, err := createContainerCmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return []metadata.Package{}, fmt.Errorf("failed to create container: %s %s", e.Stderr, e)
		}
		return []metadata.Package{}, fmt.Errorf("failed to create container: %s", err)
	}

	var packages []metadata.Package
	status, err := getStatus(string(containerId))
	if err != nil {
		return []metadata.Package{}, err
	}

	if status != nil {
		packages = parseStatDB(status)
	}

	statusD, err := getStatusD(string(containerId))
	if err != nil {
		return []metadata.Package{}, err
	}

	if statusD != nil {
		packages = append(packages, parseStatDBEntries(statusD)...)
	}

	return packages, nil
}

func parseStatDBEntries(t *tar.Reader) (packages []metadata.Package) {
	for {
		h, err := t.Next()
		if err == io.EOF {
			break // End of archive
		}

		if h.Typeflag == tar.TypeDir {
			continue
		}

		all, err := ioutil.ReadAll(t)
		if err != nil {
			continue
		}
		packageEntry, err := ParseStatDBEntry(string(all))
		if err == nil {
			packages = append(packages, packageEntry)
		}
	}
	return packages
}

func ParseStatDBEntry(content string) (metadata.Package, error) {
	pkg := metadata.Package{}

	if strings.TrimSpace(content) == "" {
		return pkg, errors.New("invalid StatDB entry")
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

func parseStatDB(r io.Reader) []metadata.Package {
	packages := make([]metadata.Package, 0)

	statDBString, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error read statDBString: %s", err)
	}

	statDBEntries := strings.Split(string(statDBString), "\n\n")
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
