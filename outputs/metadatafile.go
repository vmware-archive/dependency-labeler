package outputs

import (
	"encoding/json"
	"github.com/pivotal/deplab/metadata"
	"log"
	"os"
)

func WriteMetadataFile(md metadata.Metadata, metadataFilePath string) {
	metadataFile, err := os.OpenFile(metadataFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("no such file: %s\n", metadataFilePath)
	}
	encoder := json.NewEncoder(metadataFile)
	_ = encoder.Encode(md)
}
