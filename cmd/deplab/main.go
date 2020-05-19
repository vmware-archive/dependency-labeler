package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/pivotal/deplab/pkg/deplab"

	"github.com/spf13/cobra"
)

var (
	additionalSourceFilePaths []string
	inputImage                string
	inputImageTar             string
	outputImageTar            string
	gitPaths                  []string
	metadataFilePath          string
	dpkgFilePath              string
	tag                       string
	additionalSourceUrls      []string
	ignoreValidationErrors    bool
)

func init() {
	rootCmd.Flags().StringArrayVarP(&gitPaths, "git", "g", []string{}, "`path` to a directory under git revision control")
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "image which will be analysed by deplab. Cannot be used with --image-tar flag")
	rootCmd.Flags().StringVarP(&inputImageTar, "image-tar", "p", "", "`path` to tarball of input image. Cannot be used with --image flag")
	rootCmd.Flags().StringVarP(&outputImageTar, "output-tar", "o", "", "`path` to write a tarball of the image to")
	rootCmd.Flags().StringVarP(&metadataFilePath, "metadata-file", "m", "", "write metadata to this file at the given `path`")
	rootCmd.Flags().StringVarP(&dpkgFilePath, "dpkg-file", "d", "", "write dpkg list metadata in (modified) 'dpkg -l' format to a file at this `path`")
	rootCmd.Flags().StringVarP(&tag, "tag", "t", "", "tags the output image")
	rootCmd.Flags().StringArrayVarP(&additionalSourceUrls, "additional-source-url", "u", []string{}, "`url` to the source of an added dependency")
	rootCmd.Flags().StringArrayVarP(&additionalSourceFilePaths, "additional-sources-file", "a", []string{}, "`path` to file describing additional sources")
	rootCmd.Flags().BoolVar(&ignoreValidationErrors, "ignore-validation-errors", false, "Set flag to ignore validation errors")
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.deplab.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,
	Version: deplab.Version,

	PreRunE: validateFlags,

	Run: run,
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if !isFlagSet(cmd, "image") && !isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: requires one of --image or --image-tar")
	} else if isFlagSet(cmd, "image") && isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: cannot accept both --image and --image-tar")
	}

	if !isFlagSet(cmd, "metadata-file") && !isFlagSet(cmd, "dpkg-file") && !isFlagSet(cmd, "output-tar") {
		return fmt.Errorf("ERROR: requires one of --metadata-file, --dpkg-file, or --output-tar")
	}

	return nil
}

func isFlagSet(cmd *cobra.Command, flagName string) bool {
	flagSet := cmd.Flags()
	flag, err := flagSet.GetString(flagName)
	if err != nil {
		return false
	}

	return flag != ""
}

func run(_ *cobra.Command, _ []string) {
	err := deplab.Run(
		common.RunParams{
			InputImageTarPath:         inputImageTar,
			InputImage:                inputImage,
			GitPaths:                  gitPaths,
			Tag:                       tag,
			OutputImageTar:            outputImageTar,
			MetadataFilePath:          metadataFilePath,
			DpkgFilePath:              dpkgFilePath,
			AdditionalSourceUrls:      additionalSourceUrls,
			AdditionalSourceFilePaths: additionalSourceFilePaths,
			IgnoreValidationErrors:    ignoreValidationErrors,
		})
	if err != nil {
		log.Fatalf("deplab failed to run. %s\n", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
