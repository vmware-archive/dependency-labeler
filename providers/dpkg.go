package providers

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pivotal/deplab/docker"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/pivotal/deplab/metadata"
)

const DebianPackageListSourceType = "debian_package_list"
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
			Type: DebianPackageListSourceType,
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

	grep := exec.Command("docker", "run", "--rm",
		"--entrypoint", "grep", imageName,
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

func getDebianPackages(imageName string) ([]metadata.Package, error) {
	var packages []metadata.Package

	statusPackages, err := listPackagesFromStatus(imageName)

	if err != nil {
		return []metadata.Package{}, err
	}

	packages = append(packages, statusPackages...)

	statusDPackages, err := listPackagesFromStatusD(imageName)

	if err != nil {
		return []metadata.Package{}, err
	}

	packages = append(packages, statusDPackages...)

	collator := collate.New(language.BritishEnglish)
	sort.Slice(packages, func(i, j int) bool{
		return collator.CompareString(packages[i].Package, packages[j].Package) < 0
	})

	return packages, nil
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

func listPackagesFromStatusD(imageName string) (packages []metadata.Package, err error) {
	path := "/var/lib/dpkg/status.d"
	t, err := docker.ReadFromImage(imageName, path)

	if err != nil {
		return packages, fmt.Errorf("error retrieving %s %s", path, err)
	}

	for {
		h, err := t.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return packages, err
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
	return packages, nil
}

func listPackagesFromStatus(imageName string) (packages []metadata.Package, err error) {
	path := "/var/lib/dpkg/status"
	t, err := docker.ReadFromImage(imageName, path)

	if err != nil {
		return packages, fmt.Errorf("error retrieving %s %s", path, err)
	}

	_, err = t.Next()
	if err != nil {
		if err == io.EOF {
			return packages, nil
		}
		return packages, fmt.Errorf("error reading status: %s", err)
	}

	statDBString, err := ioutil.ReadAll(t)
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

	return packages, nil
}

func getUpstreamVersion(input string) string {
	version := strings.Split(input, "-")[0]
	if strings.Contains(version, ":") {
		version = strings.Split(version, ":")[1]
	}
	return version
}
