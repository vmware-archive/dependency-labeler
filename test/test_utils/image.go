package test_utils

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	"path/filepath"
)

type MockImage struct {
	path string
	config *v1.ConfigFile
}

func (m MockImage) GetConfig() (*v1.ConfigFile, error) {
	return m.config, nil
}

func NewMockImageWithEmptyConfig() MockImage {
	return MockImage{
		config: &v1.ConfigFile{},
	}
}

func (m MockImage) Cleanup() {
	panic("implement me")
}

func (m MockImage) GetFileContent(string) (string, error) {
	return "", nil
}

func (m MockImage) GetDirContents(string) ([]string, error) {
	return []string{}, nil
}

func (m MockImage) AbsolutePath(string) (string, error) {
	path, err := filepath.Abs(m.path)

	Expect(err).ToNot(HaveOccurred())
	return path, err
}

func NewMockImageWithPath(path string) MockImage {
	return MockImage{
		path: path,
	}
}

func (m MockImage) ExportWithMetadata(metadata.Metadata, string, string) error {
	panic("implement me")
}