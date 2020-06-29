// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package rpm

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRPM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[rpm] RPM Suite")
}
