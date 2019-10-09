package deplab

import (
	"bytes"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pivotal/deplab/docker"
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

	var originImageTarPath string
	if inputImageTarPath != ""{
		originImageTarPath = inputImageTarPath
	}	else {

		// use crane.pull to get tar ball and put it in originImageTarPath.
		dir, err := ioutil.TempDir("", "deplab-crane-")
		if err != nil {
			log.Fatalf("Could not create temp directory. %s", err)
		}
		defer os.RemoveAll(dir)

		originImageTarPath = dir + "/image.tgz"

		pulledImage, err := crane.Pull(inputImage)
		if err != nil {
			log.Fatalf("could not pull image from url: %s; Err: %s", inputImage, err)
		}
		ref, err := name.ParseReference(inputImage)
		if err != nil {
			log.Fatalf("Could not parse input image reference: %s. Err: %s", inputImage, err)
		}
		imgTag := ref.(name.Tag)
		err = crane.Save(pulledImage, imgTag.Name(), originImageTarPath)
		if err != nil {
			log.Fatalf("could not save image to path: %s; Err: %s", originImageTarPath, err)
		}
	}


	originImage := inputImage
	if originImageTarPath != "" {
		stdout, stderr, err := runCommand("docker", "load", "-i", originImageTarPath)
		if err != nil {
			log.Fatalf("could not load docker image from tar at %s; Err: %s", originImageTarPath, stderr)
		}

		imageTag := ""
		if strings.Contains(stdout.String(), "image ID") {
			imageTag = strings.TrimPrefix(stdout.String(), "Loaded image ID:")
		} else {
			imageTag = strings.TrimPrefix(stdout.String(), "Loaded image:")
		}
		originImage = strings.TrimSpace(imageTag)
	}

	gitDependencies, archiveUrls := preprocess(gitPaths, additionalSourceFilePaths)
	additionalSourceUrls = append(additionalSourceUrls, archiveUrls...)

	dependencies, err := generateDependencies(originImage, gitDependencies, additionalSourceUrls)
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

func preprocess(gitPaths, additionalSourcesFiles []string) ([]metadata.Dependency, []string){
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

func generateDependencies(imageName string, gitDependencies []metadata.Dependency, archiveUrls []string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := providers.BuildDebianDependencyMetadata(imageName)
	if err != nil {
		log.Fatalf("debian package metadata: %s", err)
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	dependencies = append(dependencies, gitDependencies...)

	for _, archiveUrl := range archiveUrls {
		archiveMetadata, err := providers.BuildArchiveDependencyMetadata(archiveUrl)

		if err != nil {
			log.Fatalf("archive metadata: %s", err)
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
