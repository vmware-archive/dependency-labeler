package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/pivotal/deplab/metadata"

	"github.com/google/go-containerregistry/pkg/crane"
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
	Long:    `prints the deplab "io.pivotal.metadata" label in the config file of an OCI compatible image tarball to stdout.  The label will be printed in json format.`,
	PreRunE: validateInspectFlags,
	Run: func(cmd *cobra.Command, args []string) {
		var img v1.Image
		var err error
		var inputPath string
		if inputImageTar != "" {
			img, err = crane.Load(inputImageTar)
			inputPath = inputImageTar
		} else {
			img, err = crane.Pull(inputImage)
			inputPath = inputImage
		}
		if err != nil {
			log.Fatalf("deplab cannot open the provided image %s: %s", inputPath, err)
		}

		cf, err := img.ConfigFile()
		if err != nil {
			log.Fatalf("deplab cannot open the Config file for %s: %s", inputPath, err)
		}
		if label, ok := cf.Config.Labels["io.pivotal.metadata"]; !ok {
			log.Fatalf("deplab cannot find the 'io.pivotal.metadata' label on the provided image: %s", inputPath)
		} else {
			err := json.Unmarshal([]byte(label), &metadata.Metadata{})
			if err != nil {
				log.Fatalf("deplab cannot parse the label on the provided image %s, label: %s: %s", inputPath, label, err)
			}
			stdOutBuffer := bytes.Buffer{}

			err = json.Indent(&stdOutBuffer, []byte(label), "", "  ")
			if err != nil {
				log.Fatalf("deplab cannot pretty print the label of the provided image %s, label: %s: %s", inputPath, label, err)
			}

			fmt.Println(stdOutBuffer.String())
			os.Exit(0)
		}
	},
}

func validateInspectFlags(cmd *cobra.Command, args []string) error {
	if !isFlagSet(cmd, "image") && !isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: requires one of --image or --image-tar")
	} else if isFlagSet(cmd, "image") && isFlagSet(cmd, "image-tar") {
		return fmt.Errorf("ERROR: cannot accept both --image and --image-tar")
	}

	return nil
}
