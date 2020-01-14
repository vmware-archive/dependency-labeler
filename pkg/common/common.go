package common

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
