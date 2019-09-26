package deplab

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/pivotal/deplab/docker"
	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/outputs"
	"github.com/pivotal/deplab/providers"
)

var (
	DeplabVersion string
)

const UnknownDeplabVersion = "0.0.0-dev"

func Run(inputImageTar string, inputImage string, gitPaths []string, tag string, outputImageTar string, metadataFilePath string, dpkgFilePath string) {
	originImage := inputImage
	if inputImageTar != "" {
		stdout, stderr, err := runCommand("docker", "load", "-i", inputImageTar)
		if err != nil {
			log.Fatalf("could not load docker image from tar: %s", stderr)
		}

		imageTag := ""
		if strings.Contains(stdout.String(), "image ID") {
			imageTag = strings.TrimPrefix(stdout.String(), "Loaded image ID:")
		} else {
			imageTag = strings.TrimPrefix(stdout.String(), "Loaded image:")
		}
		originImage = strings.TrimSpace(imageTag)
	}
	dependencies, err := generateDependencies(originImage, gitPaths)
	if err != nil {
		log.Fatalf("error generating dependencies: %s", err)
	}
	md := metadata.Metadata{Dependencies: dependencies}

	md.Base = providers.BuildOSMetadata(originImage)

	md.Provenance = []metadata.Provenance{{
		Name:    "deplab",
		Version: GetVersion(),
		URL:     "https://github.com/pivotal/deplab",
	}}

	resp, err := docker.CreateNewImage(originImage, md, tag)
	if err != nil {
		log.Fatalf("could not create new image: %s\n", err)
	}

	newID, err := docker.GetIDOfNewImage(resp)
	if err != nil {
		log.Fatalf("could not get ID of the new image: %s\n", err)
	}

	fmt.Println(newID)

	writeOutputs(md, metadataFilePath, dpkgFilePath)

	if outputImageTar != "" {
		id := newID

		if tag != "" {
			id = tag
		}

		_, stderr, err := runCommand("docker", "save", id, "-o", outputImageTar)
		if err != nil {
			log.Fatalf("could not save docker image to tar: %s", stderr)
		}
	}

}

func GetVersion() string {
	if DeplabVersion == "" {
		return UnknownDeplabVersion
	}

	return DeplabVersion
}

func generateDependencies(imageName string, pathsToGit []string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := providers.BuildDebianDependencyMetadata(imageName)
	if err != nil {
		log.Fatalf("debian package metadata: %s", err)
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	for _, gitPath := range pathsToGit {
		gitMetadata, err := providers.BuildGitDependencyMetadata(gitPath)

		if err != nil {
			log.Fatalf("git metadata: %s", err)
		}
		dependencies = append(dependencies, gitMetadata)
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
