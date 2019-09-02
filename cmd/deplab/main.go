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
	metadataFilePath string
	dpkgFilePath     string
	tag              string
)

func init() {
	rootCmd.Flags().StringVarP(&gitPath, "git", "g", "", "`path` to a directory under git revision control")
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "image which will be analysed by deplab. Cannot be used with --image-tar flag")
	rootCmd.Flags().StringVarP(&inputImageTar, "image-tar", "p", "", "`path` to tarball of input image. Cannot be used with --image flag")
	rootCmd.Flags().StringVarP(&outputImageTar, "output-tar", "o", "", "`path` to write a tarball of the image to")
	rootCmd.Flags().StringVarP(&metadataFilePath, "metadata-file", "m", "", "write metadata to this file at the given `path`")
	rootCmd.Flags().StringVarP(&dpkgFilePath, "dpkg-file", "d", "", "write dpkg list metadata in (modified) 'dpkg -l' format to a file at this `path`")
	rootCmd.Flags().StringVarP(&tag, "tag", "t", "", "tags the output image")

	_ = rootCmd.MarkFlagRequired("git")
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,
	Version: deplab.GetVersion(),

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
