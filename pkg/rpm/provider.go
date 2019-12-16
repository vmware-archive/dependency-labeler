package rpm

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pivotal/deplab/pkg/metadata"
)

const RPMSourceType = "rpm_package_list"

type ImageInterface interface {
	AbsolutePath(string) string
}

func BuildRPMMetadata(dli ImageInterface) (metadata.Dependency, error) {
	stdOutBuffer := &strings.Builder{}

	cmd := exec.Command("rpm",
		"-qa",
		"--dbpath", dli.AbsolutePath("/var/lib/rpm"),
		"--queryformat", QueryFormat(),
	)

	cmd.Stdout = stdOutBuffer

	err := cmd.Run()

	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("executing rpm: %w", err)
	}

	if strings.TrimSpace(stdOutBuffer.String()) == "" {
		return metadata.Dependency{}, errors.New("no rpm packages data found")
	}

	allPackagesDetails := strings.Split(strings.TrimSpace(stdOutBuffer.String()), "\n")

	var packages []metadata.RpmPackage

	for _, line := range allPackagesDetails {
		packages = append(packages, UnmarshalPackages(line))
	}

	return metadata.Dependency{
		Type: RPMSourceType,
		Source: metadata.Source{
			Type:    "",
			Version: nil,
			Metadata: metadata.RpmPackageListSourceMetadata{
				Packages: packages,
			},
		},
	}, nil
}
