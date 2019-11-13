package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pivotal/deplab/metadata"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"
)

func init() {
	inspectCmd.Flags().StringVarP(&inputImageTar, "image-tar", "p", "", "`path` to tarball of input image. Cannot be used with --image flag")

	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "prints the deplab label to stdout",
	Long:  `prints the deplab "io.pivotal.metadata" label in the config file of an OCI compatible image tarball to stdout.  The label will be printed in json format.`,
	Run: func(cmd *cobra.Command, args []string) {
		img, err := crane.Load(inputImageTar)
		if err != nil {
			log.Fatalf("deplab cannot open the provided image %s: %s", inputImageTar, err)
		}
		cf, err := img.ConfigFile()
		if err != nil {
			log.Fatalf("deplab cannot open the Config file for %s: %s", inputImageTar, err)
		}
		if label, ok := cf.Config.Labels["io.pivotal.metadata"]; !ok {
			log.Fatalf("deplab cannot find the 'io.pivotal.metadata' label on the provided image: %s", inputImageTar)
		} else {
			err := json.Unmarshal([]byte(label), &metadata.Metadata{})
			if err != nil {
				log.Fatalf("deplab cannot parse the label on the provided image %s, label: %s: %s", inputImageTar, label, err)
			}
			stdOutBuffer := bytes.Buffer{}

			err = json.Indent(&stdOutBuffer, []byte(label), "", "  ")
			if err != nil {
				log.Fatalf("deplab cannot pretty print the label of the provided image %s, label: %s: %s", inputImageTar, label, err)
			}

			fmt.Println(stdOutBuffer.String())
			os.Exit(0)
		}
	},
}
