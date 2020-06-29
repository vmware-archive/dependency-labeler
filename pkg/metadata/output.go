// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package metadata

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteMetadataFile(md Metadata, metadataFilePath string) error {
	metadataFile, err := os.Create(metadataFilePath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", metadataFilePath, err)
	}
	encoder := json.NewEncoder(metadataFile)
	err = encoder.Encode(md)
	if err != nil {
		return fmt.Errorf("could not write metadata file: %w", err)
	}
	return nil
}
