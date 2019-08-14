package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pivotal/deplab/metadata"
)

func BuildDebianDependencyMetadata(imageName string) (metadata.Dependency, error) {
	packages, err := getDebianPackages(imageName)

	if len(packages) != 0 {
		sources, _ := getAptSources(imageName)

		dpkgList := metadata.Dependency{
			Type: "debian_package_list",
			Source: metadata.Source{
				Type: "inline",
				Metadata: metadata.DebianPackageListSourceMetadata{
					Packages:   packages,
					AptSources: sources,
				},
			},
		}

		return dpkgList, nil
	}

	return metadata.Dependency{}, err
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

func getDebianPackages(imageName string) ([]metadata.Package, error) {
	query := `{
		"package":"${Package}",
		"version":"${Version}",
		"architecture":"${architecture}",
		"source":{
			"package":"${source:Package}",
			"version":"${source:Version}",
			"upstreamVersion":"${source:Upstream-Version}"
		}
	},`

	dpkgQuery := exec.Command("docker", "run", "--rm", imageName, "dpkg-query", "-W", "-f", query)

	out, err := dpkgQuery.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "executable file not found in $PATH") {
			log.Print("This image does not contain dpkg, so skipping dpkg dependencies.")
			return []metadata.Package{}, nil
		}
		return []metadata.Package{}, fmt.Errorf("dpkgQuery failed: %s, with error: %s\n", string(out), err.Error())
	}

	amendedOut := string(out)
	pattern := regexp.MustCompile(`{`)
	loc := pattern.FindIndex(out)
	amendedOut = amendedOut[loc[0] : len(amendedOut)-1]
	amendedOut = "[" + amendedOut + "]"

	decoder := json.NewDecoder(strings.NewReader(amendedOut))
	var packages []metadata.Package
	err = decoder.Decode(&packages)
	if err != nil {
		return []metadata.Package{}, fmt.Errorf("unable to decode pkg: %s\n", err.Error())
	}

	return packages, nil
}
