package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
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
		return "", errors.Wrapf(err, "could not encode source metadata.")
	}
	version := hex.EncodeToString(hash.Sum(nil))
	return version, nil
}
