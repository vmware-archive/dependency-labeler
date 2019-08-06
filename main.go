package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jhoonb/archivex"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("no parameters passed in. Expecting image as parameter")
	}

	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"), client.FromEnv)
	if err != nil {
		panic(err)
	}

	stdOutBuffer := bytes.Buffer{}
	inputImage := os.Args[1]

	tar := new(archivex.TarFile)
	tar.CreateWriter("docker context", &stdOutBuffer)
	tar.Add("Dockerfile", strings.NewReader("FROM "+inputImage), nil)
	tar.Close()

	opt := types.ImageBuildOptions{
		Labels: map[string]string{
			"io.pivotal.metadata": "metadata here",
		},
	}

	resp, err := cli.ImageBuild(context.Background(), &stdOutBuffer, opt)
	if err != nil {
		log.Fatalf("could not build image: %s\n", err)
	}

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
			log.Fatalln("could not find the new image ID")
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "error reading line")
			continue
		}

		if line.Error != "" {
			log.Fatalf("error building image: %s\n", line.Error)
		} else if line.Aux.ID != "" {
			fmt.Println(line.Aux.ID)
			return
		}
	}

}
