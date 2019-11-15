package deplab

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pivotal/deplab/pkg/additionalsources"

	"github.com/pivotal/deplab/pkg/git"

	"github.com/pivotal/deplab/pkg/dpkg"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/pivotal/deplab/pkg/os"

	"github.com/pivotal/deplab/pkg/image"

	"github.com/pkg/errors"
)

var Version = "0.0.0-dev"

type RunParams struct {
	InputImageTarPath         string
	InputImage                string
	GitPaths                  []string
	Tag                       string
	OutputImageTar            string
	MetadataFilePath          string
	DpkgFilePath              string
	AdditionalSourceUrls      []string
	AdditionalSourceFilePaths []string
	IgnoreValidationErrors    bool
}

func Run(params RunParams) error {
	dli, err := image.NewDeplabImage(params.InputImage, params.InputImageTarPath)
	if err != nil {
		return errors.Wrapf(err, "could not load image.")
	}
	defer dli.Cleanup()

	gitDependencies, archiveUrls, err := preprocess(params.GitPaths, params.AdditionalSourceFilePaths, params.IgnoreValidationErrors)
	if err != nil {
		return errors.Wrapf(err, "could not preprocess provided data.")
	}
	params.AdditionalSourceUrls = append(params.AdditionalSourceUrls, archiveUrls...)

	err = additionalsources.ValidateURLs(params.AdditionalSourceUrls, http.Head)
	if err != nil {
		errMsg := fmt.Sprintf("failed to validate additional source url: %s", err)
		if params.IgnoreValidationErrors {
			log.Printf("warning: %s", errMsg)
		} else {
			return errors.Errorf("error: %s", errMsg)
		}
	}

	dependencies, err := generateDependencies(dli, gitDependencies, params.AdditionalSourceUrls)
	if err != nil {
		return errors.Wrapf(err, "error generating dependencies")
	}
	md := metadata.Metadata{Dependencies: dependencies}

	md.Base = os.BuildOSMetadata(dli)

	md.Provenance = []metadata.Provenance{{
		Name:    "deplab",
		Version: Version,
		URL:     "https://github.com/pivotal/deplab",
	}}

	if params.OutputImageTar != "" {
		err = dli.ExportWithMetadata(md, params.OutputImageTar, params.Tag)

		if err != nil {
			return errors.Wrapf(err, "error exporting tar to %s", params.OutputImageTar)
		}
	}

	err = writeOutputs(md, params.MetadataFilePath, params.DpkgFilePath)
	if err != nil {
		return errors.Wrapf(err, "could not write outputs.")
	}

	return nil
}

func preprocess(gitPaths, additionalSourcesFiles []string, ignoreValidationErrors bool) ([]metadata.Dependency, []string, error) {
	var archiveUrls []string
	var gitDependencies []metadata.Dependency
	for _, additionalSourcesFile := range additionalSourcesFiles {
		archiveUrlsFromAdditionalSourcesFile, gitVcsFromAdditionalSourcesFile, err := additionalsources.ParseAdditionalSourcesFile(additionalSourcesFile)
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
		gitMetadata, err := git.BuildDependencyMetadata(gitPath)

		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not build git metadata")
		}
		gitDependencies = append(gitDependencies, gitMetadata)
	}

	return gitDependencies, archiveUrls, nil
}

func generateDependencies(dli image.Image, gitDependencies []metadata.Dependency, archiveUrls []string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := dpkg.BuildDependencyMetadata(dli)
	if err != nil {
		return dependencies, errors.Wrapf(err, "Could not generate debian package dependencies.")
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	dependencies = append(dependencies, gitDependencies...)

	for _, archiveUrl := range archiveUrls {
		archiveMetadata, err := additionalsources.BuildArchiveDependencyMetadata(archiveUrl)

		if err != nil {
			return dependencies, errors.Wrapf(err, "Could not generate archive dependency metadata.")
		}
		dependencies = append(dependencies, archiveMetadata)
	}
	return dependencies, nil
}

func writeOutputs(md metadata.Metadata, metadataFilePath, dpkgFilePath string) error {
	if metadataFilePath != "" {
		err := metadata.WriteMetadataFile(md, metadataFilePath)
		if err != nil {
			return errors.Wrapf(err, "could not write metadata file.")
		}
	}

	if dpkgFilePath != "" {
		err := dpkg.WriteDpkgFile(md, dpkgFilePath, Version)
		if err != nil {
			return errors.Wrapf(err, "could not write dpkg file.")
		}
	}

	return nil
}
