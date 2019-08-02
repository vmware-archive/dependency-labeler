package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jhoonb/archivex"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
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

	resp, _ := cli.ImageBuild(context.Background(), &stdOutBuffer, opt)

	rd := bufio.NewReader(resp.Body)

	for {
		line, err := rd.ReadString('\n')

		if err == io.EOF {
			break
		}
		if strings.HasPrefix(line, `{"aux":`) {
			out := struct {
				Aux struct {
					ID string
				}
			}{}

			_ = json.Unmarshal([]byte(line), &out)
			fmt.Println(out.Aux.ID)
			return
		}

	}
}
