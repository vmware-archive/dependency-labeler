// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package osrelease_test

import (
	"fmt"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	. "github.com/pivotal/deplab/pkg/osrelease"
)

type MockImage struct {
	fileContent string
	fileContentError error
	dirFileNames []string
}

func (m MockImage) GetConfig() (*v1.ConfigFile, error) {
	panic("implement me")
}

func (m MockImage) Cleanup() {
	panic("implement me")
}

func (m MockImage) GetFileContent(string) (string, error) {
	return m.fileContent, m.fileContentError
}

func (m MockImage) GetDirFileNames(string, bool) ([]string, error) {
	return m.dirFileNames, nil
}

func (m MockImage) GetDirContents(string) ([]string, error) {
	panic("implement me")
}

func (m MockImage) AbsolutePath(string) (string, error) {
	panic("implement me")
}

func (m MockImage) ExportWithMetadata(metadata.Metadata, string, string) error {
	panic("implement me")
}

var _ = Describe("OsRelease", func() {
	Describe("BuildOSMetadata", func() {
		Context("when the image has os-release", func() {
			Context("when os-release is malformed", func() {
				It("returns properly formatted os metadata", func() {
					fileContent := "bad\\bad"

					Expect(BuildOSMetadata(MockImage{fileContent:fileContent, fileContentError:nil})).To(Equal(metadata.UnknownBase))
				})
			})
			Context("when os-release is properly formatted", func(){
				It("returns properly formated os metadata", func() {
					fileContent := "PRETTY_NAME=\"some-pretty-name\"\nNAME=\"some-name\"\nVERSION_ID=\"some-version-id\"\nVERSION=\"some-version\"\nVERSION_CODENAME=some-version-codename\nID=some-id\nHOME_URL=\"some-url\"\n"

					expectedMetadata := metadata.Base {
						"id": "some-id",
						"pretty_name": "some-pretty-name",
						"version": "some-version",
						"home_url": "some-url",
						"version_id": "some-version-id",
						"version_codename": "some-version-codename",
						"name": "some-name",
					}

					Expect(BuildOSMetadata(MockImage{fileContent:fileContent, fileContentError:nil})).To(Equal(expectedMetadata))
				})
			})
		})

		Context("when the image does not have os-release", func() {
			Context("when the image has /bin/ash", func() {
				It("returns busybox os metadata", func() {
					dirFileNames := []string{"sh", "bash", "rash", "dash", "stash", "ash"}
					fileContentError := fmt.Errorf("could not find the file")

					Expect(BuildOSMetadata(MockImage{dirFileNames:dirFileNames, fileContentError:fileContentError})).To(Equal(metadata.BusyboxBase))
				})
			})

			Context("when the image does not have /bin/ash", func() {
				Context("when there are 3 or less file/folders in the root directory", func() {
					It("returns scratch base os metadata", func() {
						fileContentError := fmt.Errorf("could not find the file")
						rootContentNames := []string{"bin", "dev", "sys"}

						Expect(BuildOSMetadata(MockImage{dirFileNames:rootContentNames, fileContentError:fileContentError})).To(Equal(metadata.ScratchBase))
					})
				})

				It("returns unknown base os metadata", func() {
					dirFileNames := []string{"sh", "bash", "rash", "dash", "stash"}
					fileContentError := fmt.Errorf("could not find the file")

					Expect(BuildOSMetadata(MockImage{dirFileNames:dirFileNames, fileContentError:fileContentError})).To(Equal(metadata.UnknownBase))
				})
			})
		})
	})
})