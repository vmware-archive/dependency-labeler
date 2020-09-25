// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package cnb_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCnb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cnb Suite")
}
