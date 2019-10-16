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
		It("accepts url for which head function returns any 2xx status code", func() {
			statusCode := 199
			err := providers.ValidateURLs([]string{
				"http://example.com",
				"http://example.com",
				"http://example.com",
				"http://example.com",
			}, func(_ string) (*http.Response, error) {
				statusCode++
				return &http.Response{StatusCode: statusCode}, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not accept url for which head function returns a non 2xx status code", func() {
			err := providers.ValidateURLs([]string{
				"http://example.com",
			}, func(_ string) (*http.Response, error) {
				return &http.Response{StatusCode: 400}, nil
			})
			Expect(err).To(HaveOccurred())
		})

		It("does not accept url for which head function returns an error", func() {
			err := providers.ValidateURLs([]string{
				"http://example.com",
			}, func(u string) (*http.Response, error) {
				return nil, fmt.Errorf("some error with %s", u)
			})
			Expect(err).To(MatchError(ContainSubstring("http://example.com")))
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
				return &http.Response{StatusCode: 200}, nil
			})

			Expect(err).To(HaveOccurred())
		})
	})
})
