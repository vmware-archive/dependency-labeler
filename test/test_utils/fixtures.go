package test_utils

import "github.com/pivotal/deplab/pkg/metadata"

var MetadataSample = metadata.Metadata{
	Dependencies: []metadata.Dependency{
		{
			Type: "debian_package_list",
			Source: metadata.Source{
				Version: map[string]interface{}{
					"sha256": "some-sha",
				},
				Metadata: metadata.DebianPackageListSourceMetadata{
					Packages: []metadata.Package{
						{
							Package:      "foobar",
							Version:      "0.42.0-version",
							Architecture: "amd46",
							Source: metadata.PackageSource{
								Package:         "foobar",
								Version:         "0.42.0-source",
								UpstreamVersion: "0.42.0-upstream",
							},
						},
					},
					AptSources: nil,
				},
			},
		},
	},
}