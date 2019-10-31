package deplab

import (
	"fmt"
	"log"
	"net/http"

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

func Run(inputImageTarPath string, inputImage string,
	gitPaths []string, tag string, outputImageTar string,
	metadataFilePath string, dpkgFilePath string,
	additionalSourceUrls []string, additionalSourceFilePaths []string,
	ignoreValidationErrors bool) error {
	dli, err := rootfs.NewDeplabImage(inputImage, inputImageTarPath)
	if err != nil {
		return errors.Wrapf(err, "could not load image.")
	}
	defer dli.Cleanup()

	gitDependencies, archiveUrls, err := preprocess(gitPaths, additionalSourceFilePaths, ignoreValidationErrors)
	if err != nil {
		return errors.Wrapf(err, "could not preprocess provided data.")
	}
	additionalSourceUrls = append(additionalSourceUrls, archiveUrls...)

	err = providers.ValidateURLs(additionalSourceUrls, http.Head)
	if err != nil {
		errMsg := fmt.Sprintf("failed to validate additional source url: %s", err)
		if ignoreValidationErrors {
			log.Printf("warning: %s", errMsg)
		} else {
			return errors.Errorf("error: %s", errMsg)
		}
	}

	dependencies, err := generateDependencies(dli, gitDependencies, additionalSourceUrls)
	if err != nil {
		return errors.Wrapf(err, "error generating dependencies")
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
			return errors.Wrapf(err, "error exporting tar to %s", outputImageTar)
		}
	}

	err = writeOutputs(md, metadataFilePath, dpkgFilePath)
	if err != nil {
		return errors.Wrapf(err, "could not write outputs.")
	}

	return nil
}

func GetVersion() string {
	if DeplabVersion == "" {
		return UnknownDeplabVersion
	}

	return DeplabVersion
}

func preprocess(gitPaths, additionalSourcesFiles []string, ignoreValidationErrors bool) ([]metadata.Dependency, []string, error) {
	var archiveUrls []string
	var gitDependencies []metadata.Dependency
	for _, additionalSourcesFile := range additionalSourcesFiles {
		archiveUrlsFromAdditionalSourcesFile, gitVcsFromAdditionalSourcesFile, err := preprocessors.ParseAdditionalSourcesFile(additionalSourcesFile)
		if err != nil {
			errMsg := fmt.Sprintf("could not parse additional sources file: %s, %s", additionalSourcesFile, err)
			if ignoreValidationErrors {
				log.Printf("warning: %s", errMsg)
			} else {
				return nil, nil, errors.Errorf("error: %s", errMsg)
			}
		}
		archiveUrls = append(archiveUrls, archiveUrlsFromAdditionalSourcesFile...)
		gitDependencies = append(gitDependencies, gitVcsFromAdditionalSourcesFile...)
	}

	for _, gitPath := range gitPaths {
		gitMetadata, err := preprocessors.BuildGitDependencyMetadata(gitPath)

		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not build git metadata")
		}
		gitDependencies = append(gitDependencies, gitMetadata)
	}

	return gitDependencies, archiveUrls, nil
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

func writeOutputs(md metadata.Metadata, metadataFilePath, dpkgFilePath string) error {
	if metadataFilePath != "" {
		err := outputs.WriteMetadataFile(md, metadataFilePath)
		if err != nil {
			return errors.Wrapf(err, "could not write metadata file.")
		}
	}

	if dpkgFilePath != "" {
		err := outputs.WriteDpkgFile(md, dpkgFilePath, GetVersion())
		if err != nil {
			return errors.Wrapf(err, "could not write dpkg file.")
		}
	}

	return nil
}
