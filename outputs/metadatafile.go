package outputs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pivotal/deplab/metadata"
)

func WriteMetadataFile(md metadata.Metadata, metadataFilePath string) error {
	metadataFile, err := os.Create(metadataFilePath)
	if err != nil {
		return fmt.Errorf("no such file: %s\n", metadataFilePath)
	}
	encoder := json.NewEncoder(metadataFile)
	_ = encoder.Encode(md)
	return nil
}
