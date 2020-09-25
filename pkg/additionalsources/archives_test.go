// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package additionalsources_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/vmware-tanzu/dependency-labeler/pkg/additionalsources"
)

var _ = Describe("additionalsources", func() {
	Describe("ValidateURLs", func() {
		It("accepts url for which head function returns any 2xx status code", func() {
			statusCode := 199
			err := ValidateURLs([]string{
				"http://example.com/file.zip",
				"http://example.com/file.zip",
				"http://example.com/file.zip",
				"http://example.com/file.zip",
			}, func(_ string) (*http.Response, error) {
				statusCode++
				return &http.Response{StatusCode: statusCode}, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not accept url for which head function returns a non 2xx status code", func() {
			err := ValidateURLs([]string{
				"http://example.com/file.zip",
			}, func(_ string) (*http.Response, error) {
				return &http.Response{StatusCode: 400}, nil
			})
			Expect(err).To(HaveOccurred())
		})

		It("does not accept url for which head function returns an error", func() {
			err := ValidateURLs([]string{
				"http://example.com/file.zip",
			}, func(u string) (*http.Response, error) {
				return nil, fmt.Errorf("some error with %s", u)
			})
			Expect(err).To(MatchError(ContainSubstring("http://example.com/file.zip")))
		})

		It("returns an error if any of the url is not reachable", func() {
			err := ValidateURLs([]string{
				"http://example.com/file.zip",
				"http://example.com/this-is-a-404/file.zip",
				"http://example.com/foo/bar/file.zip",
				"http://example.com/file.zip",
			}, func(u string) (*http.Response, error) {
				if u == "http://example.com/foo/bar/file.zip" {
					return nil, fmt.Errorf("some error with %s", u)
				}
				if u == "http://example.com/this-is-a-404/file.zip" {
					return &http.Response{StatusCode: 404}, nil
				}
				return &http.Response{StatusCode: 200}, nil
			})

			Expect(err).To(HaveOccurred())
		})

		It("returns an error if any of the url does not have a valid extension", func() {
			err := ValidateURLs([]string{
				"http://example.com/file.zip",
				"http://example.com/file.zap",
			}, func(u string) (*http.Response, error) {
				return &http.Response{StatusCode: 200}, nil
			})

			Expect(err).To(MatchError(ContainSubstring("http://example.com/file.zap")))
		})

		It("returns an error if there is an invalid file extension", func() {
			err := ValidateURLs([]string{"http://www.somewebsite.com/file_wrong.ext"}, func(u string) (*http.Response, error) {
				return &http.Response{StatusCode: 200}, nil
			})

			Expect(err).To(HaveOccurred())
		})

		DescribeTable("with a valid extension", func(validateUrl string) {
			err := ValidateURLs([]string{validateUrl}, func(u string) (*http.Response, error) {
				return &http.Response{StatusCode: 200}, nil
			})

			Expect(err).ToNot(HaveOccurred())
		},
			//https://en.wikipedia.org/wiki/Tar_(computing)#Suffixes_for_compressed_files
			generateEntries(
				"7z",
				"tar.bz2",
				"tar.gz",
				"tar.lz",
				"tar.lzma",
				"tar.lzo",
				"tar.xz",
				"tar.Z",
				"tar.zst",
				"taz",
				"taZ",
				"tb2",
				"tbz",
				"tbz2",
				"tgz",
				"tlz",
				"tpz",
				"txz",
				"tZ",
				"tz2",
				"tzst",
				"tar.bz2",
				"zip",
				"zip#with-an-ignored-fragment",
			)...,
		)
	})
})

func generateEntries(extensions ...string) []TableEntry {
	var entries []TableEntry

	for _, e := range extensions {
		entries = append(entries, Entry("of "+e, "http://www.somewebsite.com/file."+e))
	}
	return entries
}
