package rpm

import (
	"github.com/pivotal/deplab/pkg/metadata"
	"os/exec"
	"strings"
)



const RPMSourceType = "rpm_package_list"

type  ImageInterface interface{
	AbsolutePath(string) string
}


func BuildRPMMetadata(dli ImageInterface) (metadata.Dependency, error) {
	stdOutBuffer := &strings.Builder{}

	cmd := exec.Command("rpm",
		"-qa",
		"--dbpath", dli.AbsolutePath("/var/lib/rpm"),
		"--queryformat", "%{NAME}\n",
		)

	cmd.Stdout = stdOutBuffer

	_ = cmd.Run()

	packagesNames := strings.Split(strings.TrimSpace(stdOutBuffer.String()), "\n")
	var packages []metadata.Package
	for _, pkg := range packagesNames{
		packages = append(packages, metadata.Package{Package: pkg})
	}

	return metadata.Dependency{
		Type:   RPMSourceType,
		Source: metadata.Source{
			Type:     "",
			Version:  nil,
			Metadata: metadata.RpmPackageListSourceMetadata{
				Packages: packages,
			},
		},
	}, nil
}
