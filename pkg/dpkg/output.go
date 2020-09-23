// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package dpkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"
)

func WriteDpkgFile(md metadata.Metadata, dpkgFilePath string, deplabVersion string) error {
	f, err := os.Create(dpkgFilePath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", dpkgFilePath, err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("could not close file: %s; error: %s\n", dpkgFilePath, err)
		}
	}()

	dep, err := findDpkgListInMetadata(md)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	pkgs := dep.Source.Metadata.(metadata.DebianPackageListSourceMetadata).Packages

	tHeader := []string{"||/", "Name", "Version", "Architecture", "Description"}
	tRows := make([][]string, 0)
	tMaxLen := []int{3, 4, 7, 12, 11}

	for _, pkg := range pkgs {
		tRow := []string{"ii", pkg.Package, pkg.Version, pkg.Architecture, "Description intentionally left blank"}
		for i, v := range tMaxLen {
			if utf8.RuneCountInString(tRow[i]) > v {
				tMaxLen[i] = utf8.RuneCountInString(tRow[i])
			}
		}
		tRows = append(tRows, tRow)
	}

	header := `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
`

	sha, ok := dep.Source.Version["sha256"].(string)
	if !ok {
		return fmt.Errorf("version in dpkg list was not a string")
	}

	unpaddedFmtString := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds\n", tMaxLen[0]+1, tMaxLen[1]+1, tMaxLen[2]+1, tMaxLen[3]+1, tMaxLen[4]+1)
	fmtString := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n", tMaxLen[0], tMaxLen[1], tMaxLen[2], tMaxLen[3], tMaxLen[4])

	df := dpkgFile{w: f}

	df.
		printf("deplab SHASUM: %s\n", sha).
		printf("deplab version: %s\n\n", deplabVersion).
		printf(header).
		printf(fmtString, tHeader[0], tHeader[1], tHeader[2], tHeader[3], tHeader[4]).
		printf(unpaddedFmtString,
			"+++-",
			fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[1])),
			fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[2])),
			fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[3])),
			fmt.Sprintf("%s", strings.Repeat("=", tMaxLen[4])),
		)

	for _, v := range tRows {
		df.printf(fmtString, v[0], v[1], v[2], v[3], v[4])
	}

	if df.err != nil {
		return fmt.Errorf("Could not write to file: %s\n", dpkgFilePath)
	}
	return nil
}

type dpkgFile struct {
	err error
	w   io.Writer
}

func (df *dpkgFile) printf(format string, a ...interface{}) *dpkgFile {
	if df.err == nil {
		_, err := fmt.Fprintf(df.w, format, a...)
		df.err = err
	}
	return df
}

func findDpkgListInMetadata(md metadata.Metadata) (metadata.Dependency, error) {
	for _, v := range md.Dependencies {
		if v.Type == metadata.DebianPackageListSourceType {
			return v, nil
		}
	}

	return metadata.Dependency{}, fmt.Errorf("could not find %s in metadata", metadata.DebianPackageListSourceType)
}
