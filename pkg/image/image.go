// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package image

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"

	"github.com/pivotal/deplab/pkg/metadata"

	"github.com/containerd/containerd/reference/docker"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

type Image interface {
	GetFileContent(string) (string, error)
	GetDirContents(string) ([]string, error)
	GetDirFileNames(string, bool) ([]string, error)
	AbsolutePath(string) (string, error)
	GetConfig() (*v1.ConfigFile, error)
	ExportWithMetadata(metadata.Metadata, string, string) error
}

type ExportableImage interface {
	ExportWithMetadata(metadata.Metadata, string, string) error
	Cleanup()
}

type RootFSImage struct {
	rootFS RootFS
	image  v1.Image
}

func (dli RootFSImage) GetConfig() (*v1.ConfigFile, error) {
	return dli.image.ConfigFile()
}

func NewDeplabImage(inputImage, inputImageTarPath string) (RootFSImage, error) {
	var (
		image v1.Image
		err   error
	)

	if inputImage != "" {
		image, err = crane.Pull(inputImage)
		if err != nil {
			return RootFSImage{}, fmt.Errorf("failed to pull %s: %w", inputImage, err)
		}
	} else if inputImageTarPath != "" {
		image, err = crane.Load(inputImageTarPath)
		if err != nil {
			return RootFSImage{}, fmt.Errorf("failed to load %s: %w", inputImageTarPath, err)
		}
	} else {
		return RootFSImage{}, fmt.Errorf("you must provide either an inputImage or inputImageTarPath parameter")
	}

	// this folder is unnecessary and may contain folders with bad permissions
	rootFS, err := NewRootFS(image, []string{"usr/share/doc/"})
	if err != nil {
		return RootFSImage{}, fmt.Errorf("could not create new image: %w", err)
	}

	return RootFSImage{image: image, rootFS: rootFS}, nil
}

func (dli *RootFSImage) Cleanup() {
	dli.rootFS.Cleanup()
}

func (dli RootFSImage) ExportWithMetadata(metadata metadata.Metadata, path string, tag string) error {
	err := dli.setMetadata(metadata)
	if err != nil {
		return fmt.Errorf("error setting metadata: %w", err)
	}

	err = dli.export(path, tag)
	if err != nil {
		return fmt.Errorf("error exporting tar to %s: %w", path, err)
	}
	return nil
}

func (dli RootFSImage) GetFileContent(s string) (string, error) {
	return dli.rootFS.GetFileContent(s)
}

func (dli RootFSImage) GetDirContents(s string) ([]string, error) {
	return dli.rootFS.GetDirContents(s)
}

func (dli RootFSImage) GetDirFileNames(s string, i bool) ([]string, error) {
	return dli.rootFS.GetDirFileNames(s, i)
}

func (dli *RootFSImage) setMetadata(metadata metadata.Metadata) error {
	config, err := dli.image.ConfigFile()
	if err != nil {
		return fmt.Errorf("could not find config file in image: %w", err)
	}
	md, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}
	if config.Config.Labels == nil {
		config.Config.Labels = map[string]string{}
	}

	config.Config.Labels["io.deplab.metadata"] = string(md)

	dli.image, err = mutate.Config(dli.image, config.Config)
	if err != nil {
		return fmt.Errorf("could not mutate config in image: %w", err)
	}

	return nil
}

func (dli *RootFSImage) export(path string, tag string) error {
	var actualTag string
	if tag == "" {
		h, err := dli.image.Digest()
		if err != nil {
			return fmt.Errorf("could not retrieve image digest: %w", err)
		}
		actualTag = h.String()
	} else {
		actualTag = tag
	}

	name, err := docker.ParseDockerRef(actualTag)
	if err != nil {
		return fmt.Errorf("tag %s is invalid: %w", actualTag, err)
	}

	err = crane.Save(dli.image, name.String(), path)
	if err != nil {
		return fmt.Errorf("could not export to %s: %w", path, err)
	}

	return nil
}

func (dli RootFSImage) AbsolutePath(absPath string) (string, error) {
	joinedPath := path.Join(dli.rootFS.rootfsLocation, absPath)
	patheee, err := filepath.Abs(joinedPath)
	if err != nil {
		return "", fmt.Errorf("could not create image absolute path, %s: %w", absPath, err)
	}
	return patheee, nil
}
