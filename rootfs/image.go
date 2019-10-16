package rootfs

import (
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/pivotal/deplab/metadata"
	"github.com/pkg/errors"
)

type Image struct {
	rootfs RootFS
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

	rootfs, _ := New(image)

	return Image{image: image, rootfs: rootfs}, nil
}

func (dli *Image) Cleanup() {
	dli.rootfs.Cleanup()
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
	return dli.rootfs.GetFileContent(s)
}

func (dli *Image) GetDirContents(s string) ([]string, error) {
	return dli.rootfs.GetDirContents(s)
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
		h, _ := dli.image.Digest()
		actualTag = h.String()
	} else {
		actualTag = tag
	}

	err := crane.Save(dli.image, actualTag, path)

	return errors.Wrapf(err, "could not export to %s: %s", path, err)
}
