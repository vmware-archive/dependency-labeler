package main

import (
	"fmt"
	"os"

	"github.com/pivotal/deplab"
	"github.com/spf13/cobra"
)

var (
	inputImage       string
	inputImageTar    string
	outputImageTar   string
	gitPaths         []string
	metadataFilePath string
	dpkgFilePath     string
	tag              string
)

func init() {
	rootCmd.Flags().StringArrayVarP(&gitPaths, "git", "g", []string{}, "`path` to a directory under git revision control")
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

	PreRunE: preRunE,

	Run: run,
}

func preRunE(cmd *cobra.Command, args []string) error {
	flagset := cmd.Flags()
	img, err := flagset.GetString("image")
	if err != nil {
		return fmt.Errorf("error processing flag: %s", err)
	}
	imgTar, err := flagset.GetString("image-tar")
	if err != nil {
		return fmt.Errorf("error processing flag: %s", err)
	}

	if img == "" && imgTar == "" {
		return fmt.Errorf("ERROR: requires one of --image or --image-tar")
	} else if img != "" && imgTar != "" {
		return fmt.Errorf("ERROR: cannot accept both --image and --image-tar")
	}
	return nil
}

func run(cmd *cobra.Command, args []string) {
	deplab.Run(inputImageTar, inputImage, gitPaths, tag, outputImageTar, metadataFilePath, dpkgFilePath)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
