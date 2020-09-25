// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package metadata_test

import (
	"io/ioutil"

	"github.com/vmware-tanzu/dependency-labeler/test/test_utils"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

var _ = Describe("outputs", func() {
	Describe("WriteMetadataFile", func() {
		DescribeTable("when the file can be written", func(path string) {
			defer test_utils.CleanupFile(path)

			err := WriteMetadataFile(test_utils.MetadataSample, path)
			Expect(err).ToNot(HaveOccurred())

			_, err = ioutil.ReadFile(path)

			Expect(err).ToNot(HaveOccurred())
		},
			Entry("the file exists", test_utils.ExistingFileName()),
			Entry("the file does not exists", test_utils.NonExistingFileName()),
		)

		Describe("and metadata can't be written", func() {
			It("returns an error", func() {
				err := WriteMetadataFile(test_utils.MetadataSample, "a-path-that-does-not-exist/foo.dpkg")

				Expect(err).To(MatchError(ContainSubstring("a-path-that-does-not-exist/foo.dpkg")))
			})
		})

		Describe("when the file already has content", func() {
			It("wipes the original file", func() {
				By("creating a metadata file")
				filePath := test_utils.ExistingFileName()
				err := WriteMetadataFile(test_utils.MetadataSample, filePath)
				Expect(err).ToNot(HaveOccurred())

				originalContent := test_utils.AppendContent(filePath)

				By("running it again against the same file")
				err = WriteMetadataFile(test_utils.MetadataSample, filePath)
				Expect(err).ToNot(HaveOccurred())

				newBytes, err := ioutil.ReadFile(filePath)
				Expect(string(newBytes), err).To(Equal(string(originalContent)))
			})
		})
	})
})
