package deplab

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/pivotal/deplab/rootfs"

	"github.com/pkg/errors"

	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/outputs"
	"github.com/pivotal/deplab/preprocessors"
	"github.com/pivotal/deplab/providers"
)

var (
	DeplabVersion string
)

const UnknownDeplabVersion = "0.0.0-dev"

func Run(inputImageTarPath string, inputImage string, gitPaths []string, tag string, outputImageTar string, metadataFilePath string, dpkgFilePath string, additionalSourceUrls []string, additionalSourceFilePaths []string) {
	dli, err := rootfs.NewDeplabImage(inputImage, inputImageTarPath)
	if err != nil {
		log.Fatalf("could not load image: %s", err)
	}
	defer dli.Cleanup()

	gitDependencies, archiveUrls := preprocess(gitPaths, additionalSourceFilePaths)
	additionalSourceUrls = append(additionalSourceUrls, archiveUrls...)

	err = providers.ValidateURLs(additionalSourceUrls, http.Head)
	if err != nil {
		log.Fatalf("error validating additional source url: %s", err)
	}

	dependencies, err := generateDependencies(dli, gitDependencies, additionalSourceUrls)
	if err != nil {
		log.Fatalf("error generating dependencies: %s", err)
	}
	md := metadata.Metadata{Dependencies: dependencies}

	md.Base = providers.BuildOSMetadata(dli)

	md.Provenance = []metadata.Provenance{{
		Name:    "deplab",
		Version: GetVersion(),
		URL:     "https://github.com/pivotal/deplab",
	}}

	if outputImageTar != "" {
		err = dli.ExportWithMetadata(md, outputImageTar, tag)

		if err != nil {
			log.Fatalf("error exporting tar to %s: %s", outputImageTar, err)
		}
	}

	writeOutputs(md, metadataFilePath, dpkgFilePath)
}

func GetVersion() string {
	if DeplabVersion == "" {
		return UnknownDeplabVersion
	}

	return DeplabVersion
}

func preprocess(gitPaths, additionalSourcesFiles []string) ([]metadata.Dependency, []string) {
	var archiveUrls []string
	var gitDependencies []metadata.Dependency
	for _, additionalSourcesFile := range additionalSourcesFiles {
		archiveUrlsFromAdditionalSourcesFile, gitVcsFromAdditionalSourcesFile, err := preprocessors.ParseAdditionalSourcesFile(additionalSourcesFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("could not parse additional sources file: %s", additionalSourcesFile)))
		}
		archiveUrls = append(archiveUrls, archiveUrlsFromAdditionalSourcesFile...)
		gitDependencies = append(gitDependencies, gitVcsFromAdditionalSourcesFile...)
	}

	for _, gitPath := range gitPaths {
		gitMetadata, err := preprocessors.BuildGitDependencyMetadata(gitPath)

		if err != nil {
			log.Fatalf("git metadata: %s", err)
		}
		gitDependencies = append(gitDependencies, gitMetadata)
	}

	return gitDependencies, archiveUrls
}

func generateDependencies(dli rootfs.Image, gitDependencies []metadata.Dependency, archiveUrls []string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := providers.BuildDebianDependencyMetadata(dli)
	if err != nil {
		return dependencies, errors.Wrapf(err, "Could not generate debian package dependencies.")
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	dependencies = append(dependencies, gitDependencies...)

	for _, archiveUrl := range archiveUrls {
		archiveMetadata, err := providers.BuildArchiveDependencyMetadata(archiveUrl)

		if err != nil {
			return dependencies, errors.Wrapf(err, "Could not generate archive dependency metadata.")
		}
		dependencies = append(dependencies, archiveMetadata)
	}
	return dependencies, nil
}

func writeOutputs(md metadata.Metadata, metadataFilePath, dpkgFilePath string) {
	if metadataFilePath != "" {
		outputs.WriteMetadataFile(md, metadataFilePath)
	}

	if dpkgFilePath != "" {
		outputs.WriteDpkgFile(md, dpkgFilePath, GetVersion())
	}
}

func runCommand(cmd string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	dockerLoad := exec.Command(cmd, args...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	dockerLoad.Stdout = stdout
	dockerLoad.Stderr = stderr

	err := dockerLoad.Run()
	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}
