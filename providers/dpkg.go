package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pivotal/deplab/pkg/metadata"
)

func BuildDebianDependencyMetadata(imageName string) (metadata.Dependency, error) {
	packages, err := getDebianPackages(imageName)

	if len(packages) != 0 {
		dpkgList := metadata.Dependency{
			Type: "debian_package_list",
			Source: metadata.Source{
				Type: "inline",
				Metadata: metadata.DebianPackageListSourceMetadata{
					Packages: packages,
				},
			},
		}

		return dpkgList, nil
	}

	return metadata.Dependency{}, err
}

func getDebianPackages(imageName string) ([]metadata.Package, error) {
	query := "{\"package\":\"${Package}\", \"version\":\"${Version}\", \"architecture\":\"${architecture}\", \"source\":{\"package\":\"${source:Package}\", \"version\":\"${source:Version}\", \"upstreamVersion\":\"${source:Upstream-Version}\"}},"

	dpkgQuery := exec.Command("docker", "run", "--rm", imageName, "dpkg-query", "-W", "-f="+query)

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
