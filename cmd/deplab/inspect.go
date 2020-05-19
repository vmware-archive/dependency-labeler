package main

import (
	"fmt"

	"github.com/pivotal/deplab/pkg/deplab"

	"github.com/spf13/cobra"
)

func init() {
	inspectCmd.Flags().StringVarP(&inputImageTar, "image-tar", "p", "", "`path` to tarball of input image. Cannot be used with --image flag")
	inspectCmd.Flags().StringVarP(&inputImage, "image", "i", "", "image which will be inspected by deplab. Cannot be used with --image-tar flag")

	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:     "inspect",
	Short:   "prints the deplab label to stdout",
	Long:    `prints the deplab "io.deplab.metadata" label in the config file of an OCI compatible image to stdout.  The label will be printed in json format.`,
	PreRunE: validateInspectFlags,
	RunE: func(_ *cobra.Command, _ []string) error {
		return deplab.RunInspect(inputImage, inputImageTar)
	},
}

func validateInspectFlags(cmd *cobra.Command, _ []string) error {
	if !isFlagSet(cmd, "image") && !isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: requires one of --image or --image-tar")
	} else if isFlagSet(cmd, "image") && isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: cannot accept both --image and --image-tar")
	}

	return nil
}
