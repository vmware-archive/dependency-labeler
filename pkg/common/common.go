// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type RunParams struct {
	InputImageTarPath         string
	InputImage                string
	GitPaths                  []string
	Tag                       string
	OutputImageTar            string
	MetadataFilePath          string
	DpkgFilePath              string
	AdditionalSourceUrls      []string
	AdditionalSourceFilePaths []string
	IgnoreValidationErrors    bool
}

func Digest(sourceMetadata interface{}) (string, error) {
	hash := sha256.New()
	encoder := json.NewEncoder(hash)
	err := encoder.Encode(sourceMetadata)
	if err != nil {
		return "", fmt.Errorf("could not encode source metadata: %w", err)
	}
	version := hex.EncodeToString(hash.Sum(nil))
	return version, nil
}
