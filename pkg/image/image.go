package image

import (
	"encoding/json"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/pkg/errors"
)

type Image struct {
	rootFS RootFS
	image  v1.Image
}

func NewDeplabImage(inputImage, inputImageTarPath string) (Image, error) {
	var (
		image v1.Image
		err   error
	)

	if inputImage != "" {
		image, err = crane.Pull(inputImage)
		if err != nil {
			return Image{}, errors.Wrapf(err, "failed to pull %s: %s", inputImage, err)
		}
	} else if inputImageTarPath != "" {
		image, err = crane.Load(inputImageTarPath)
		if err != nil {
			return Image{}, errors.Wrapf(err, "failed to load %s: %s", inputImageTarPath, err)
		}
	} else {
		return Image{}, errors.New("You must provide either an inputImage or inputImageTarPath parameter")
	}

	rootFS, err := NewRootFS(image)
	if err != nil {
		return Image{}, errors.Wrapf(err, "could not create new image")
	}

	return Image{image: image, rootFS: rootFS}, nil
}

func (dli *Image) Cleanup() {
	dli.rootFS.Cleanup()
}

func (dli *Image) ExportWithMetadata(metadata metadata.Metadata, path string, tag string) error {
	err := dli.setMetadata(metadata)
	if err != nil {
		return errors.Wrapf(err, "error setting metadata: %s", err)
	}

	err = dli.export(path, tag)
	if err != nil {
		return errors.Wrapf(err, "error exporting tar to %s: %s", path, err)
	}
	return nil
}

func (dli *Image) GetFileContent(s string) (string, error) {
	return dli.rootFS.GetFileContent(s)
}

func (dli *Image) GetDirContents(s string) ([]string, error) {
	return dli.rootFS.GetDirContents(s)
}

func (dli *Image) setMetadata(metadata metadata.Metadata) error {
	config, err := dli.image.ConfigFile()
	if err != nil {
		return errors.Wrapf(err, "could not find config file in image : %s", err)
	}
	md, err := json.Marshal(metadata)
	if err != nil {
		return errors.Wrapf(err, "could not marshal json : %s", err)
	}
	if config.Config.Labels == nil {
		config.Config.Labels = map[string]string{}
	}

	config.Config.Labels["io.pivotal.metadata"] = string(md)

	dli.image, err = mutate.Config(dli.image, config.Config)
	if err != nil {
		return errors.Wrapf(err, "could not mutate config in image : %s", err)
	}

	return nil
}

func (dli *Image) export(path string, tag string) error {
	var actualTag string
	if tag == "" {
		h, err := dli.image.Digest()
		if err != nil {
			return errors.Wrapf(err, "could not retrieve image digest")
		}
		actualTag = h.String()
	} else {
		actualTag = tag
	}

	err := crane.Save(dli.image, actualTag, path)

	return errors.Wrapf(err, "could not export to %s: %s", path, err)
}
