package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pivotal/deplab/builder"
	"github.com/pivotal/deplab/providers"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/spf13/cobra"
)

var inputImage string
var gitPath string

var deplabVersion string

func init() {
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "Image for the metadata to be added to")
	rootCmd.Flags().StringVarP(&gitPath, "git", "g", "", "Path to directory under git revision control")
	err := rootCmd.MarkFlagRequired("image")
	if err != nil {
		log.Fatalf("image name is required\n")
	}
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,
	Version: deplabVersion,

	Run: func(cmd *cobra.Command, args []string) {
		if !builder.IsValidImageName(inputImage) {
			log.Fatalf("invalid image name: %s\n", inputImage)
		}

		dependencies, err := GenerateDependencies(inputImage, gitPath)
		if err != nil {
			log.Fatalf("error generating dependencies: %s", err)
		}
		md := metadata.Metadata{Dependencies: dependencies}

		osMetadata, err := providers.BuildOSMetadata(inputImage)
		md.Base = osMetadata

		resp, err := builder.CreateNewImage(inputImage, md)
		if err != nil {
			log.Fatalf("could not create new image: %s\n", err)
		}

		newID, err := builder.GetIDOfNewImage(resp)
		if err != nil {
			log.Fatalf("could not get ID of the new image: %s\n", err)
		}
		fmt.Println(newID)
	},
}

func main() {
	if deplabVersion == "" {
		rootCmd.Version = "0.0.0-dev"
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GenerateDependencies(imageName, pathToGit string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := providers.BuildDebianDependencyMetadata(imageName)
	if err != nil {
		log.Fatalf("debian package metadata: %s", err)
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	if gitPath != "" {
		gitMetadata, err := providers.BuildGitDependencyMetadata(pathToGit)
		if err != nil {
			log.Fatalf("git metadata: %s", err)
		}
		dependencies = append(dependencies, gitMetadata)
	}

	return dependencies, nil
}
