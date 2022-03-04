// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

module github.com/vmware-tanzu/dependency-labeler

go 1.16

require (
	github.com/containerd/containerd v1.6.1
	github.com/docker/docker v20.10.12+incompatible
	github.com/google/go-containerregistry v0.8.0
	github.com/joho/godotenv v1.4.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1
	github.com/spf13/cobra v1.3.0
	golang.org/x/text v0.3.7
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/moby/sys/mount v0.3.1 // indirect
