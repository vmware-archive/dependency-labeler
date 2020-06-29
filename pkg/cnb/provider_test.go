// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package cnb_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal/deplab/pkg/cnb"
	"github.com/pivotal/deplab/pkg/common"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/test/test_utils"
)

var _ = Describe("Cnb", func() {
	Describe("Provider", func() {
		Context("when the image has no cnb label", func() {
			It("does not modify the metadata content", func() {
				Expect(Provider(test_utils.NewMockImageWithEmptyConfig(), common.RunParams{}, metadata.Metadata{})).To(Equal(metadata.Metadata{}))
			})
		})
	})
})
