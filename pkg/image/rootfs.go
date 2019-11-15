package image

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/google/go-containerregistry/pkg/v1/mutate"

	"github.com/docker/docker/pkg/archive"

	"github.com/pkg/errors"
)

type Interface interface {
	GetFileContent(path string) (string, error)
	GetDirContents(path string) ([]string, error)
}

type RootFS struct {
	rootfsLocation string
}

func (rfs *RootFS) GetDirContents(path string) ([]string, error) {
	var fileContents []string
	files, err := ioutil.ReadDir(filepath.Join(rfs.rootfsLocation, path))
	if err != nil {
		return fileContents, errors.Wrapf(err, "could not find directory in rootFS: %s", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		thisFileContent, err := rfs.GetFileContent(filepath.Join(path, f.Name()))
		if err != nil {
			return fileContents, errors.Wrapf(err, "could not find file in directory in rootFS: %s", err)
		}
		fileContents = append(fileContents, thisFileContent)
	}

	return fileContents, nil
}

func (rfs *RootFS) GetFileContent(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(rfs.rootfsLocation, path))
	if err != nil {
		return "", errors.Wrapf(err, "could not find file in rootFS: %s", err)
	}
	return string(fileBytes), nil
}

func NewRootFS(image v1.Image) (RootFS, error) {
	var (
		rootFS string
		err    error
	)

	f, err := ioutil.TempFile("", "image")
	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not create temp file.")
	}

	fs := mutate.Extract(image)

	rootFS, err = ioutil.TempDir("", "deplab-rootFS-")
	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not create rootFS temp directory.")
	}

	err = archive.Untar(fs, rootFS, &archive.TarOptions{NoLchown: true})
	if err != nil {
		return RootFS{}, errors.Wrapf(err, "Could not untar from tar %s to temp directory %s.", f.Name(), rootFS)
	}

	return RootFS{rootfsLocation: rootFS}, nil
}

func (rfs *RootFS) Cleanup() {
	err := os.RemoveAll(rfs.rootfsLocation)
	if err != nil {
		log.Printf("could not clean up rootFS location: %s. %s\n", rfs.rootfsLocation, err)
	}
}
