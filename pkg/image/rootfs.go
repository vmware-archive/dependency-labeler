package image

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"github.com/docker/docker/pkg/archive"
)

type Interface interface {
	GetFileContent(path string) (string, error)
	GetDirContents(path string) ([]string, error)
}

type RootFS struct {
	rootfsLocation string
}

const RootfsPrefix = "deplab-rootFS-"

func (rfs *RootFS) GetDirContents(path string) ([]string, error) {
	var fileContents []string
	files, err := ioutil.ReadDir(filepath.Join(rfs.rootfsLocation, path))
	if err != nil {
		return fileContents, fmt.Errorf("could not find directory in rootFS: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		thisFileContent, err := rfs.GetFileContent(filepath.Join(path, f.Name()))
		if err != nil {
			return fileContents, fmt.Errorf("could not find file in directory in rootFS: %w", err)
		}
		fileContents = append(fileContents, thisFileContent)
	}

	return fileContents, nil
}

func (rfs *RootFS) GetFileContent(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(rfs.rootfsLocation, path))
	if err != nil {
		return "", fmt.Errorf("could not find file in rootFS: %w", err)
	}
	return string(fileBytes), nil
}

func NewRootFS(image v1.Image, excludePatterns []string) (RootFS, error) {
	var (
		rootFS string
		err    error
	)

	f, err := ioutil.TempFile("", "image")
	if err != nil {
		return RootFS{}, fmt.Errorf("could not create temp file: %w", err)
	}

	fs := mutate.Extract(image)

	rootFS, err = ioutil.TempDir("", RootfsPrefix)
	if err != nil {
		return RootFS{}, fmt.Errorf("could not create rootFS temp directory: %w", err)
	}

	err = archive.Untar(fs, rootFS, &archive.TarOptions{
		ExcludePatterns: excludePatterns,
		NoLchown:        true,
		InUserNS:        true,
	})

	if err != nil {
		return RootFS{}, fmt.Errorf("could not untar from tar %s to temp directory %s: %w", f.Name(), rootFS, err)
	}

	return RootFS{rootfsLocation: rootFS}, nil
}

func (rfs *RootFS) Cleanup() {
	err := os.RemoveAll(rfs.rootfsLocation)
	if err != nil {
		log.Printf("could not clean up rootFS location: %s. %s\n", rfs.rootfsLocation, err)
	}
}
