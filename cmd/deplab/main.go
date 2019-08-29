package main

import (
	"fmt"
	"github.com/pivotal/deplab"
	"github.com/spf13/cobra"
	"log"
	"os"
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

	PreRun: preRun,

	Run: run,
}

func preRun(cmd *cobra.Command, args []string) {
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
}

func run(cmd *cobra.Command, args []string) {
	deplab.Run(inputImageTar, inputImage, gitPath, tag, outputImageTar, metadataFilePath, dpkgFilePath)
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
