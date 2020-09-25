// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package cnb_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmware-tanzu/dependency-labeler/pkg/cnb"
	"github.com/vmware-tanzu/dependency-labeler/pkg/common"
	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
	"github.com/vmware-tanzu/dependency-labeler/test/test_utils"
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
