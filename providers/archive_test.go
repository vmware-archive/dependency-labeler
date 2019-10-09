package providers_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/providers"
)

var _ = Describe("Providers/Archive", func() {
	Describe("ValidateURLs", func() {
		It("Accepts valid url for which head function does not return an error", func() {
			err := providers.ValidateURLs([]string{
				"http://example.com",
			}, func(_ string) (*http.Response, error) {
				return nil, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if any of the url is not reachable", func() {
			err := providers.ValidateURLs([]string{
				"http://example.com",
				"http://example.com/this-is-a-404",
				"http://example.com/foo/bar",
				"http://example.com",
			}, func(u string) (*http.Response, error) {
				if u == "http://example.com/foo/bar" {
					return nil, fmt.Errorf("some error with %s", u)
				}
				if u == "http://example.com/this-is-a-404" {
					return &http.Response{StatusCode: 404}, nil
				}
				return nil, nil
			})

			Expect(err).To(MatchError(
				SatisfyAll(
					ContainSubstring("some error"),
					ContainSubstring("http://example.com/foo/bar"),
				)))
		})
	})
})
