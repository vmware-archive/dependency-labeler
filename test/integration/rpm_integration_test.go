// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package integration_test

import (
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
)

var _ = Describe("[rpm] deplab rpm", func() {
	Context("with an image with an rpm database", func() {
		Context("rpm cli installed on the PATH", func() {
			It("returns rpm metadata", func() {
				metadataLabel := runDeplabAgainstTar(
					getTestAssetPath("image-archives/photon.tgz"))

				rpmPackages := selectRpmDependencies(metadataLabel.Dependencies)
				Expect(rpmPackages).To(HaveLen(1))
				rpmPackage := rpmPackages[0]
				Expect(rpmPackage.Type).To(Equal(metadata.RPMPackageListSourceType))
				Expect(rpmPackage.Source.Type).To(Equal("inline"))
				packages := rpmPackage.Source.Metadata.(map[string]interface{})["packages"].([]interface{})
				Expect(packages).To(HaveLen(34))
				Expect(ArePackagesSorted(packages)).To(BeTrue())

				By("generating a sha256 digest of the metadata content as version")
				Expect(rpmPackage.Source.Version["sha256"]).To(MatchRegexp(`^[0-9a-f]{64}$`))
			})
		})
		Context("rpm cli is not installed on the PATH", func() {
			It("provides a helpful error message", func() {
				PATH := os.Getenv("PATH")
				Expect(os.Setenv("PATH", "")).ToNot(HaveOccurred())

				defer func() {
					Expect(os.Setenv("PATH", PATH)).ToNot(HaveOccurred())
				}()

				f, err := ioutil.TempFile("", "")
				Expect(err).ToNot(HaveOccurred())
				defer func() {
					Expect(os.Remove(f.Name())).ToNot(HaveOccurred())
				}()

				By("executing it")
				args := []string{"--image-tar", getTestAssetPath("image-archives/photon.tgz"), "--git", pathToGitRepo, "--metadata-file", f.Name()}
				_, stdErr := runDepLab(args, 1)

				errorOutput := strings.TrimSpace(string(getContentsOfReader(stdErr)))

				Expect(errorOutput).To(SatisfyAll(
					ContainSubstring("an rpm database exists at"),
					ContainSubstring("but rpm is not installed and available on your path")))
			})
		})
	})
	Context("image without rpm database", func() {
		It("returns rpm metadata", func() {
			metadataLabel := runDeplabAgainstTar(
				getTestAssetPath("image-archives/tiny.tgz"))

			rpmPackages := selectRpmDependencies(metadataLabel.Dependencies)
			Expect(rpmPackages).To(BeEmpty())
		})
	})
})

func selectRpmDependencies(dependencies []metadata.Dependency) []metadata.Dependency {
	var rpmDependencies []metadata.Dependency
	for _, dependency := range dependencies {
		if dependency.Type == metadata.RPMPackageListSourceType {
			rpmDependencies = append(rpmDependencies, dependency)
		}
	}
	return rpmDependencies
}
