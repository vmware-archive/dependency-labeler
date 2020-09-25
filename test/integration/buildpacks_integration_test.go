// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package integration_test

import (
	"encoding/json"
	"sort"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Describe("inspect", func() {
		Context("with an image with buildpack metadata.toml file", func() {
			It("can find a buildpack dependency", func() {
				stdOut, _ := runDepLab([]string{
					"inspect",
					"--image-tar", getTestAssetPath("image-archives/scratch-with-buildpack-metadata.tgz"),
				}, 0)

				md := metadata.Metadata{}
				err := json.NewDecoder(stdOut).Decode(&md)

				Expect(err).ToNot(HaveOccurred())

				buildpacksDependency, ok := filterBuildpacksDependencies(md.Dependencies)
				Expect(ok).To(BeTrue())

				Expect(buildpacksDependency.Type).To(Equal(metadata.BuildpackMetadataType))

				By("Parsing buildpacks and bill of materials")
				buildpacksSourceMetadata := buildpacksDependency.Source.Metadata.(map[string]interface{})
				Expect(buildpacksSourceMetadata["buildpacks"].([]interface{})).To(
					ConsistOf(map[string]interface{}{
						"id":      "org.cloudfoundry.openjdk",
						"version": "v1.0.86",
					}),
				)

				Expect(buildpacksSourceMetadata["bom"].([]interface{})).To(
					ConsistOf(SatisfyAll(
						HaveKeyWithValue("name", "openjdk-jdk"),
						HaveKeyWithValue("version", "11.0.6"),
					)),
				)

				Expect(buildpacksSourceMetadata["launcher"]).ToNot(BeEmpty())

				By("ordering the buildpacks to ensure reproducibility")
				buildpacks := buildpacksSourceMetadata["buildpacks"].([]interface{})
				Expect(areBuildpacksSorted(buildpacks)).To(BeTrue())
				By("ordering the bill of materials to ensure reproducibility")
				boms := buildpacksSourceMetadata["bom"].([]interface{})
				Expect(areBomsSorted(boms)).To(BeTrue())

				By("generating a sha256 digest of the metadata content as version")
				Expect(buildpacksDependency.Source.Version["sha256"]).To(MatchRegexp(`^[0-9a-f]{64}$`))
			})
		})
	})

	Describe("root", func() {
		Context("with an image with buildpack metadata.toml file", func() {
			It("can find a buildpack dependency", func() {
				md := runDeplabAgainstTar(getTestAssetPath("image-archives/scratch-with-buildpack-metadata.tgz"))

				buildpacksDependency, ok := filterBuildpacksDependencies(md.Dependencies)
				Expect(ok).To(BeTrue())

				Expect(buildpacksDependency.Type).To(Equal(metadata.BuildpackMetadataType))

				By("Parsing buildpacks and bill of materials")
				buildpacksSourceMetadata := buildpacksDependency.Source.Metadata.(map[string]interface{})
				Expect(buildpacksSourceMetadata["buildpacks"].([]interface{})).To(
					ConsistOf(map[string]interface{}{
						"id":      "org.cloudfoundry.openjdk",
						"version": "v1.0.86",
					}),
				)

				Expect(buildpacksSourceMetadata["bom"].([]interface{})).To(
					ConsistOf(SatisfyAll(
						HaveKeyWithValue("name", "openjdk-jdk"),
						HaveKeyWithValue("version", "11.0.6"),
					)),
				)

				Expect(buildpacksSourceMetadata["launcher"]).ToNot(BeEmpty())

				By("ordering the buildpacks to ensure reproducibility")
				buildpacks := buildpacksSourceMetadata["buildpacks"].([]interface{})
				Expect(areBuildpacksSorted(buildpacks)).To(BeTrue())
				By("ordering the bill of materials to ensure reproducibility")
				boms := buildpacksSourceMetadata["bom"].([]interface{})
				Expect(areBomsSorted(boms)).To(BeTrue())

				By("generating a sha256 digest of the metadata content as version")
				Expect(buildpacksDependency.Source.Version["sha256"]).To(MatchRegexp(`^[0-9a-f]{64}$`))
			})
		})
	})
})

func filterBuildpacksDependencies(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	for _, dependency := range dependencies {
		if dependency.Type == metadata.BuildpackMetadataType {
			return dependency, true
		}
	}
	return metadata.Dependency{}, false
}

func areBuildpacksSorted(buildpacks []interface{}) bool {
	collator := collate.New(language.BritishEnglish)
	return sort.SliceIsSorted(buildpacks, func(p, q int) bool {
		lhs := buildpacks[p].(map[string]interface{})
		rhs := buildpacks[q].(map[string]interface{})
		return collator.CompareString(lhs["id"].(string), rhs["id"].(string)) <= 0
	})
}

func areBomsSorted(boms []interface{}) bool {
	collator := collate.New(language.BritishEnglish)
	return sort.SliceIsSorted(boms, func(p, q int) bool {
		lhs := boms[p].(map[string]interface{})
		rhs := boms[q].(map[string]interface{})
		return collator.CompareString(lhs["name"].(string), rhs["name"].(string)) <= 0
	})
}
