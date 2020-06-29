// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	Context("with a distroless base image", func() {
		It("[remote-image] labels the image", func() {
			metadataLabel := runDeplabAgainstImage("gcr.io/distroless/base")

			Expect(metadataLabel.Base).ToNot(BeEmpty())
		})
	})
})
