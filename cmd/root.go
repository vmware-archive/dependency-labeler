package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/jhoonb/archivex"

	"github.com/spf13/cobra"
)

var inputImage string

const ValidImageNameRE = `^([a-z0-9](?:/?(?:[._-])?(?:[a-z0-9]))*)(:[a-z0-9]+(?:[._-][a-z0-9]+)*)?$`

func init() {
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "Image for the metadata to be added to")
	rootCmd.MarkFlagRequired("image")
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,

	Run: func(cmd *cobra.Command, args []string) {
		if !regexp.MustCompile(ValidImageNameRE).MatchString(inputImage) {
			log.Fatalf("invalid image name: %s\n", inputImage)
		}

		resp, err := CreateNewImage(inputImage)
		if err != nil {
			log.Fatalf("could not build image: %s\n", err)
		}

		newID, err := GetIDOfNewImage(resp)
		if err != nil {
			log.Fatalf("could not get ID from the image: %s\n", err)
		}
		fmt.Println(newID)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CreateNewImage(inputImage string) (resp types.ImageBuildResponse, err error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"), client.FromEnv)
	if err != nil {
		return resp, err
	}

	stdOutBuffer := bytes.Buffer{}

	tar := new(archivex.TarFile)
	tar.CreateWriter("docker context", &stdOutBuffer)
	tar.Add("Dockerfile", strings.NewReader("FROM "+inputImage), nil)
	tar.Close()

	opt := types.ImageBuildOptions{
		Labels: map[string]string{
			"io.pivotal.metadata": "metadata here",
		},
	}

	resp, err = cli.ImageBuild(context.Background(), &stdOutBuffer, opt)
	return resp, err
}

func GetIDOfNewImage(resp types.ImageBuildResponse) (string, error) {
	rd := json.NewDecoder(resp.Body)

	for {
		line := struct {
			Aux struct {
				ID string
			}
			Stream string
			Error  string
		}{}

		err := rd.Decode(&line)
		if err == io.EOF {
			return "", fmt.Errorf("could not find the new image ID")
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "error reading line")
			continue
		}

		if line.Error != "" {
			return "", fmt.Errorf("error building image: %s\n", line.Error)
		} else if line.Aux.ID != "" {
			return line.Aux.ID, nil
		}
	}
}
