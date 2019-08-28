package main

import (
	"bytes"
	"fmt"
	"github.com/pivotal/deplab/outputs"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pivotal/deplab/builder"
	"github.com/pivotal/deplab/providers"

	"github.com/pivotal/deplab/metadata"

	"github.com/spf13/cobra"
)

var (
	inputImage       string
	inputImageTar    string
	outputImageTar   string
	gitPath          string
	deplabVersion    string
	metadataFilePath string
	dpkgFilePath     string
	tag              string
)

func init() {
	rootCmd.Flags().StringVarP(&gitPath, "git", "g", "", "Path to directory under git revision control")
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "Image for the metadata to be added to")
	rootCmd.Flags().StringVarP(&inputImageTar, "image-tar", "p", "", "Path to tarball of input image")
	rootCmd.Flags().StringVarP(&outputImageTar, "output-tar", "o", "", "Path to write a tarball of the image to")
	rootCmd.Flags().StringVarP(&metadataFilePath, "metadata-file", "m", "", "Write metadata to this file")
	rootCmd.Flags().StringVarP(&dpkgFilePath, "dpkg-file", "d", "", "Write dpkg list metadata in (modified) `dpkg -l` format to this file")
	rootCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tags the output image")

	_ = rootCmd.MarkFlagRequired("git")
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,
	Version: deplabVersion,

	PreRun: func(cmd *cobra.Command, args []string) {
		flagset := cmd.Flags()
		img, err := flagset.GetString("image")
		if err != nil {
			log.Fatalf("error processing flag: %s", err)
		}
		imgTar, err := flagset.GetString("image-tar")
		if err != nil {
			log.Fatalf("error processing flag: %s", err)
		}

		if img == "" && imgTar == "" {
			log.Println("ERROR: requires one of --image or --image-tar")
			cmd.Usage()
			os.Exit(1)
		} else if img != "" && imgTar != "" {
			log.Println("ERROR: cannot accept both --image and --image-tar")
			cmd.Usage()
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if inputImageTar != "" {
			stdout, stderr, err := runCommand("docker", "load", "-i", inputImageTar)
			if err != nil {
				log.Fatalf("could not load docker image from tar: %s", stderr)
			}

			imageTag := strings.Trim(stdout.String(), "Loaded image: \n")
			inputImage = imageTag
		}

		dependencies, err := generateDependencies(inputImage, gitPath)
		if err != nil {
			log.Fatalf("error generating dependencies: %s", err)
		}
		md := metadata.Metadata{Dependencies: dependencies}

		md.Base = providers.BuildOSMetadata(inputImage)

		resp, err := builder.CreateNewImage(inputImage, md, tag)
		if err != nil {
			log.Fatalf("could not create new image: %s\n", err)
		}

		newID, err := builder.GetIDOfNewImage(resp)
		if err != nil {
			log.Fatalf("could not get ID of the new image: %s\n", err)
		}

		fmt.Println(newID)

		writeOutputs(md)

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
	},
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

func writeOutputs(md metadata.Metadata) {
	if metadataFilePath != "" {
		outputs.WriteMetadataFile(md, metadataFilePath)
	}

	if dpkgFilePath != "" {
		outputs.WriteDpkgFile(md, dpkgFilePath)
	}
}
