package rootfs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/mholt/archiver"
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
		return fileContents, errors.Wrapf(err, "could not find directory in rootfs: %s", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		thisFileContent, err := rfs.GetFileContent(filepath.Join(path, f.Name()))
		if err != nil {
			return fileContents, errors.Wrapf(err, "could not find file in directory in rootfs: %s", err)
		}
		fileContents = append(fileContents, thisFileContent)
	}

	return fileContents, nil
}

func (rfs *RootFS) GetFileContent(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(rfs.rootfsLocation, path))
	if err != nil {
		return "", errors.Wrapf(err, "could not find file in rootfs: %s", err)
	}
	return string(fileBytes), nil
}

func New(pathToTar string) (RootFS, error) {
	var rootfs = ""
	var err error

	image, err := crane.Load(pathToTar)
	if err != nil {
		return RootFS{}, errors.Wrapf(err, "Could not load image from path: %s", pathToTar)
	}

	f, err := ioutil.TempFile("", "image")
	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not create temp file.")
	}
	err = crane.Export(image, f)
	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not export rootfs.")
	}

	rootfs, err = ioutil.TempDir("", "deplab-rootfs-")
	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not create rootfs temp directory.")
	}

	err = archiver.NewTar().Unarchive(f.Name(), rootfs)

	if err != nil {
		return RootFS{}, errors.Wrap(err, "Could not untar to temp directory.")
	}

	return RootFS{rootfsLocation: rootfs}, nil
}

func (rfs *RootFS) Cleanup() {
	_ = os.RemoveAll(rfs.rootfsLocation)
}
