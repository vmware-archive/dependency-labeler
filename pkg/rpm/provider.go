package rpm

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/image"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/pkg/errors"

	"github.com/pivotal/deplab/pkg/metadata"
)

const RPMDbPath = "/var/lib/rpm"

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	dependency, err := BuildDependencyMetadata(dli)
	if err != nil {
		return metadata.Metadata{}, err
	}
	md.Dependencies = append(md.Dependencies, dependency)
	return md, nil
}

func BuildDependencyMetadata(dli image.Image) (metadata.Dependency, error) {

	absPath, err := dli.AbsolutePath(RPMDbPath)
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("absolute path for rpm database: %w", err)
	}

	exists, err := exists(path.Join(absPath, "Packages"))
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("rpm could not find existance of path: %w", err)
	}
	if !exists {
		return metadata.Dependency{}, nil
	}

	if !isRPMInstalled() {
		return metadata.Dependency{}, fmt.Errorf("an rpm database exists at %s but rpm is not installed and available on your path: %w", RPMDbPath, err)
	}

	query := QueryFormat()
	cmd := exec.Command("rpm",
		"-qa",
		"--dbpath", absPath,
		"--queryformat", query,
	)
	stdOutBuffer := &strings.Builder{}
	cmd.Stdout = stdOutBuffer

	err = cmd.Run()

	if err != nil {
		return metadata.Dependency{},
			fmt.Errorf("failed to execute rpm at path, %s, with query, %s: %w", absPath, query, err)
	}

	if strings.TrimSpace(stdOutBuffer.String()) == "" {
		return metadata.Dependency{}, errors.New("no rpm packages data found")
	}

	allPackagesDetails := strings.Split(strings.TrimSpace(stdOutBuffer.String()), "\n")

	var packages []metadata.RpmPackage

	for _, line := range allPackagesDetails {
		packages = append(packages, UnmarshalPackage(line))
	}
	collator := collate.New(language.BritishEnglish)
	sort.Slice(packages, func(i, j int) bool {
		return collator.CompareString(packages[i].Package, packages[j].Package) < 0
	})

	sourceMetadata := metadata.RpmPackageListSourceMetadata{
		Packages: packages,
	}

	version, err := Digest(sourceMetadata)
	if err != nil {
		return metadata.Dependency{}, errors.Wrapf(err, "Could not get digest for source metadata")
	}

	return metadata.Dependency{
		Type: metadata.RPMPackageListSourceType,
		Source: metadata.Source{
			Type: "inline",
			Version: map[string]interface{}{
				"sha256": version,
			},
			Metadata: sourceMetadata,
		},
	}, nil
}

func isRPMInstalled() bool {
	stdOutBuffer := &strings.Builder{}
	cmd := exec.Command("rpm",
		"--version",
	)

	cmd.Stdout = stdOutBuffer

	err := cmd.Run()

	return err == nil && strings.Contains(stdOutBuffer.String(), "RPM version")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Digest(sourceMetadata metadata.RpmPackageListSourceMetadata) (string, error) {
	hash := sha256.New()
	encoder := json.NewEncoder(hash)
	err := encoder.Encode(sourceMetadata)
	if err != nil {
		return "", errors.Wrapf(err, "could not encode source metadata.")
	}
	version := hex.EncodeToString(hash.Sum(nil))
	return version, nil
}
