package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/pivotal/deplab/pkg/metadata"
)

func CreateNewImage(inputImage string, md metadata.Metadata) (resp types.ImageBuildResponse, err error) {
	dockerCli, err := client.NewClientWithOpts(client.WithVersion("1.39"), client.FromEnv)
	if err != nil {
		return resp, err
	}

	dockerfileBuffer, err := createDockerFileBuffer(inputImage)
	if err != nil {
		return resp, err
	}

	mdMarshalled, err := json.Marshal(md)
	if err != nil {
		return resp, err
	}

	opt := types.ImageBuildOptions{
		Labels: map[string]string{
			"io.pivotal.metadata": string(mdMarshalled),
		},
	}

	resp, err = dockerCli.ImageBuild(context.Background(), &dockerfileBuffer, opt)
	return resp, err
}

func createDockerFileBuffer(inputImage string) (bytes.Buffer, error) {
	dockerfileBuffer := bytes.Buffer{}

	tar := new(archivex.TarFile)
	err := tar.CreateWriter("docker context", &dockerfileBuffer)
	if err != nil {
		return dockerfileBuffer, fmt.Errorf("error creating tar writer: %s\n", err.Error())
	}
	err = tar.Add("Dockerfile", strings.NewReader("FROM "+inputImage), nil)
	if err != nil {
		return dockerfileBuffer, fmt.Errorf("error adding to the tar: %s\n", err.Error())
	}
	err = tar.Close()
	if err != nil {
		return dockerfileBuffer, fmt.Errorf("error closing the tar: %s\n", err.Error())
	}

	return dockerfileBuffer, nil
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
			_, _ = fmt.Fprintln(os.Stderr, "error reading line")
			continue
		}

		if line.Error != "" {
			return "", fmt.Errorf("error building image: %s\n", line.Error)
		} else if line.Aux.ID != "" {
			return line.Aux.ID, nil
		}
	}
}
