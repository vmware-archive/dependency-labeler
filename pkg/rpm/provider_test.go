// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package rpm_test

import (
	"github.com/pivotal/deplab/pkg/common"
	"github.com/pivotal/deplab/test/test_utils"
	"io/ioutil"
	"os"

	v1 "github.com/google/go-containerregistry/pkg/v1"

	"github.com/onsi/gomega/gstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/pkg/rpm"

	"path/filepath"
)

type MockImage struct {
	path string
}

func (m MockImage) GetConfig() (*v1.ConfigFile, error) {
	panic("implement me")
}

func (m MockImage) Cleanup() {
	panic("implement me")
}

func (m MockImage) GetFileContent(string) (string, error) {
	panic("implement me")
}

func (m MockImage) GetDirFileNames(string, bool) ([]string, error) {
	panic("implement me")
}

func (m MockImage) GetDirContents(string) ([]string, error) {
	panic("implement me")
}

func (m MockImage) AbsolutePath(string) (string, error) {
	path, err := filepath.Abs(m.path)

	Expect(err).ToNot(HaveOccurred())
	return path, err
}

func (m MockImage) ExportWithMetadata(metadata.Metadata, string, string) error {
	panic("implement me")
}

var _ = Describe("Pkg/Rpm/Provider", func() {

	//rpm leaves __db.001 etc. files in the folder when it runs; we should try to clean those up
	AfterEach(func() {
		files, err := filepath.Glob("../../test/integration/assets/rpm/__db.*")
		Expect(err).ToNot(HaveOccurred())
		for _, f := range files {
			err := os.Remove(f)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("should generate list of dependencies", func() {
		md, err := rpm.Provider(MockImage{"../../test/integration/assets/rpm"}, common.RunParams{}, metadata.Metadata{})

		Expect(err).ToNot(HaveOccurred())
		packages := md.Dependencies[0].Source.Metadata.(metadata.RpmPackageListSourceMetadata).Packages
		Expect(packages).To(HaveLen(34))

		for _, p := range packages {
			Expect(p.Package).ToNot(BeEmpty())
			Expect(p).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Package":      Not(BeEmpty()),
				"Version":      Not(BeEmpty()),
				"License":      Not(BeEmpty()),
				"Architecture": Not(BeEmpty()),
				"SourceRpm":    Not(BeEmpty()),
			}))
		}
	})

	It("does not modify the metadata if no rpm database folder is found", func() {
		tempDirPath := "/tmp/this-path-does-not-exists"
		defer func() {
			_ = os.Remove(tempDirPath)
		}()
		packages, err := rpm.Provider(MockImage{tempDirPath}, common.RunParams{}, metadata.Metadata{Dependencies: []metadata.Dependency{{
			Type: "Do not touch this one!!!!!",
		}}})
		Expect(err).NotTo(HaveOccurred())

		Expect(packages).To(Equal(metadata.Metadata{Dependencies: []metadata.Dependency{{
			Type: "Do not touch this one!!!!!",
		}}}))

	})

	It("returns nil if no Package file is found in the rpm database folder", func() {
		tempDirPath, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = os.Remove(tempDirPath)
		}()
		packages, err := rpm.Provider(test_utils.NewMockImageWithPath(tempDirPath), common.RunParams{}, metadata.Metadata{
			Dependencies: []metadata.Dependency{{
				Type: "Do not touch this one!!!!!",
			}},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(packages).To(Equal(metadata.Metadata{Dependencies: []metadata.Dependency{{
			Type: "Do not touch this one!!!!!",
		}}}))
	})

	It("returns an error if rpm is not in the PATH", func() {
		PATH := os.Getenv("PATH")
		Expect(os.Setenv("PATH", "")).ToNot(HaveOccurred())

		defer func() {
			Expect(os.Setenv("PATH", PATH)).ToNot(HaveOccurred())
		}()

		_, err := rpm.Provider(MockImage{"../../test/integration/assets/rpm"}, common.RunParams{}, metadata.Metadata{})

		Expect(err).To(MatchError(SatisfyAll(
			ContainSubstring("an rpm database exists at"),
			ContainSubstring("but rpm is not installed and available on your path"))))

	})
})
