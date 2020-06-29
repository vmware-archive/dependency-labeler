// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package kpack_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKpack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kpack Suite")
}
