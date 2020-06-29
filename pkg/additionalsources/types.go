// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package additionalsources

type AdditionalSources struct {
	Archives []AdditionalSourceArchive `yml:"archives"`
	Vcs      []AdditionalSourceVcs     `yml:"vcs"`
}

type AdditionalSourceArchive struct {
	Url string `yml:"url"`
}

type AdditionalSourceVcs struct {
	Protocol string `yml:"protocol"`
	Version  string `yml:"version"`
	Url      string `yml:"url"`
}
