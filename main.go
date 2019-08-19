package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pivotal/deplab/builder"
	"github.com/pivotal/deplab/providers"

	"github.com/pivotal/deplab/metadata"

	"github.com/spf13/cobra"
)

var (
	inputImage       string
	gitPath          string
	deplabVersion    string
	metadataFilePath string
)

func init() {
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "Image for the metadata to be added to")
	rootCmd.Flags().StringVarP(&gitPath, "git", "g", "", "Path to directory under git revision control")
	rootCmd.Flags().StringVarP(&metadataFilePath, "metadata-file", "m", "", "Write metadata to this file")

	_ = rootCmd.MarkFlagRequired("image")
	_ = rootCmd.MarkFlagRequired("git")
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,
	Version: deplabVersion,

	Run: func(cmd *cobra.Command, args []string) {
		dependencies, err := generateDependencies(inputImage, gitPath)
		if err != nil {
			log.Fatalf("error generating dependencies: %s", err)
		}
		md := metadata.Metadata{Dependencies: dependencies}

		md.Base = providers.BuildOSMetadata(inputImage)

		resp, err := builder.CreateNewImage(inputImage, md)
		if err != nil {
			log.Fatalf("could not create new image: %s\n", err)
		}

		newID, err := builder.GetIDOfNewImage(resp)
		if err != nil {
			log.Fatalf("could not get ID of the new image: %s\n", err)
		}

		fmt.Println(newID)

		if metadataFilePath != "" {
			metadataFile, err := os.OpenFile(metadataFilePath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("no such file: %s\n", metadataFilePath)
			}
			encoder := json.NewEncoder(metadataFile)
			_ = encoder.Encode(md)
		}
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

func generateDependencies(imageName, pathToGit string) ([]metadata.Dependency, error) {
	var dependencies []metadata.Dependency

	dpkgList, err := providers.BuildDebianDependencyMetadata(imageName)
	if err != nil {
		log.Fatalf("debian package metadata: %s", err)
	}
	if dpkgList.Type != "" {
		dependencies = append(dependencies, dpkgList)
	}

	gitMetadata, err := providers.BuildGitDependencyMetadata(pathToGit)
	if err != nil {
		log.Fatalf("git metadata: %s", err)
	}
	dependencies = append(dependencies, gitMetadata)

	return dependencies, nil
}
