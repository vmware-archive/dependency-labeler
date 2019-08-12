package deplab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4"

	"github.com/pivotal/deplab/pkg/metadata"

	docker "github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/jhoonb/archivex"

	"github.com/spf13/cobra"
)

var inputImage string
var gitPath string

const ValidImageNameRE = `^([a-z0-9](?:/?(?:[._-])?(?:[a-z0-9]))*)(:[a-z0-9]+(?:[._-][a-z0-9]+)*)?$`

func init() {
	rootCmd.Flags().StringVarP(&inputImage, "image", "i", "", "Image for the metadata to be added to")
	rootCmd.Flags().StringVarP(&gitPath, "git", "g", "", "Path to directory under git revision control")
	err := rootCmd.MarkFlagRequired("image")
	if err != nil {
		log.Fatalf("image name is required\n")
	}
}

var rootCmd = &cobra.Command{
	Use:   "deplab",
	Short: "dependency labeler adds a metadata label to a container image",
	Long: `Dependency labeler adds information about a container image to that image's config. 
	The information can be found in a "io.pivotal.metadata" label on the output image. 
	Complete documentation is available at http://github.com/pivotal/deplab`,

	Run: func(cmd *cobra.Command, args []string) {
		if !IsValidImageName() {
			log.Fatalf("invalid image name: %s\n", inputImage)
		}

		resp, err := CreateNewImage()
		if err != nil {
			log.Fatalf("could not create new image: %s\n", err)
		}

		newID, err := GetIDOfNewImage(resp)
		if err != nil {
			log.Fatalf("could not get ID of the new image: %s\n", err)
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

func IsValidImageName() bool {
	return regexp.MustCompile(ValidImageNameRE).MatchString(inputImage)
}

func CreateNewImage() (resp types.ImageBuildResponse, err error) {
	dockerCli, err := docker.NewClientWithOpts(docker.WithVersion("1.39"), docker.FromEnv)
	if err != nil {
		return resp, err
	}

	dockerfileBuffer, err := createDockerFileBuffer()
	if err != nil {
		return resp, err
	}
	dependencies, err := GenerateDependencies(inputImage, gitPath)
	if err != nil {
		return resp, err
	}
	md := metadata.Metadata{Dependencies: dependencies}

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

func GenerateDependencies(imageName, pathToGit string) ([]metadata.Dependency, error) {

	dpkgList := buildDebianDependencyMetadata(imageName)
	dependencies := []metadata.Dependency{dpkgList}

	if gitPath != "" {
		dependencies = append(dependencies, buildGitDependencyMetadata(pathToGit))
	}

	return dependencies, nil
}

func buildDebianDependencyMetadata(imageName string) metadata.Dependency {
	packages := getDebianPackages(imageName)

	dpkgList := metadata.Dependency{
		Type: "debian_package_list",
		Source: metadata.Source{
			Type: "inline",
			Metadata: metadata.DebianPackageListSourceMetadata{
				Packages: packages,
			},
		},
	}

	return dpkgList
}

func buildGitDependencyMetadata(pathToGit string) metadata.Dependency {
	repo, err := git.PlainOpen(pathToGit)
	if err != nil {
		log.Fatalf("cannot open git repository: %s\n", err)
	}

	ref, err := repo.Head()
	if err != nil {
		log.Fatalf("cannot find head of git repository: %s\n", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		log.Fatalf("cannot find remotes for repository: %s\n", err)
	}

	refs := []string{}
	tags, _ := repo.Tags()
	tags.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			refs = append(refs, ref.Name().Short())
		}

		return nil
	})

	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: "git",
			Version: map[string]interface{}{
				"commit": ref.Hash().String(),
			},
			Metadata: metadata.GitSourceMetadata{
				URI:  remotes[0].Config().URLs[0],
				Refs: refs,
			},
		},
	}
}

func getDebianPackages(imageName string) []metadata.Package {
	query := "{\"package\":\"${Package}\", \"version\":\"${Version}\", \"architecture\":\"${architecture}\", \"source\":{\"package\":\"${source:Package}\", \"version\":\"${source:Version}\", \"upstreamVersion\":\"${source:Upstream-Version}\"}},"

	dpkgQuery := exec.Command("docker", "run", "--rm", imageName, "dpkg-query", "-W", "-f="+query)

	out, err := dpkgQuery.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "executable file not found in $PATH") {
			log.Print("This image does not contain dpkg, so skipping dpkg dependencies.")
			return []metadata.Package{}
		}
		log.Fatalf("dpkgQuery failed: %s, with error: %s\n", string(out), err.Error())
	}

	amendedOut := string(out)
	pattern := regexp.MustCompile(`{`)
	loc := pattern.FindIndex(out)
	amendedOut = amendedOut[loc[0] : len(amendedOut)-1]
	amendedOut = "[" + amendedOut + "]"

	decoder := json.NewDecoder(strings.NewReader(amendedOut))
	var packages []metadata.Package
	err = decoder.Decode(&packages)
	if err != nil {
		log.Fatalf("unable to decode pkg: %s\n", err.Error())
	}

	return packages
}

func createDockerFileBuffer() (bytes.Buffer, error) {
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
