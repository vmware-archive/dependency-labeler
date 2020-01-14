package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/test/test_utils"
)

var _ = Describe("Merge", func() {
	Context("provenance on both original and current", func() {
		It("concatenates provenances from both", func() {
			originalProvenance := metadata.Provenance{
				Name:    "original",
				Version: "original-v1",
				URL:     "http://example.org/original",
			}
			labelMetadata := metadata.Metadata{
				Provenance: []metadata.Provenance{originalProvenance},
			}

			currentProvenance := metadata.Provenance{
				Name:    "current",
				Version: "current-v1",
				URL:     "http://example.org/original",
			}
			current := metadata.Metadata{
				Provenance: []metadata.Provenance{currentProvenance},
			}

			result, warnings := metadata.Merge(labelMetadata, current)
			Expect(result.Provenance).To(ConsistOf(originalProvenance, currentProvenance))
			Expect(warnings).To(BeEmpty())
		})
	})

	Context("base on both original and current", func() {
		Context("when original and current match", func() {
			It("retains only the base from the current metadata", func() {
				originalBase := metadata.Base{
					"name":    "original",
					"version": "original-version",
				}
				original := metadata.Metadata{
					Base: originalBase,
				}

				current := metadata.Metadata{
					Base: originalBase,
				}

				result, warnings := metadata.Merge(original, current)
				Expect(result.Base).To(Equal(originalBase))
				Expect(warnings).To(BeEmpty())
			})
		})

		Context("when original and current don't match", func() {
			It("retains only the base from the current metadata and emits a warning", func() {
				originalBase := metadata.Base{
					"name":    "original",
					"version": "original-version",
				}
				original := metadata.Metadata{
					Base: originalBase,
				}
				currentBase := metadata.Base{
					"name": "current",
				}
				current := metadata.Metadata{
					Base: currentBase,
				}

				result, warnings := metadata.Merge(original, current)
				Expect(result.Base).To(Equal(currentBase))
				Expect(warnings).To(ConsistOf(metadata.Warning("base")))
			})
		})
	})

	Context("git dependencies on original", func() {
		It("retains the git dependencies from the original metadata", func() {
			originalGit := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.GitSourceType,
				},
			}
			original := metadata.Metadata{
				Dependencies: []metadata.Dependency{
					originalGit,
					originalGit,
				},
			}

			current := metadata.Metadata{
				Dependencies: []metadata.Dependency{},
			}

			result, warnings := metadata.Merge(original, current)
			Expect(result.Dependencies).To(Equal([]metadata.Dependency{originalGit, originalGit}))
			Expect(warnings).To(BeEmpty())
		})
	})

	Context("dpkg list dependencies on both original and current", func() {
		Context("when original and current match", func() {
			It("retains only the dpkg list dependencies from the current metadata", func() {
				originalDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "some-sha256",
						},
					},
				}

				currentDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "some-sha256",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalDpkg},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentDpkg},
				})

				Expect(warnings).To(BeEmpty())

				dpkg, ok := test_utils.SelectDpkgDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(dpkg).To(Equal(currentDpkg))
			})
		})

		Context("when original and current don't match", func() {
			It("retains only the dpkg list dependencies from the current metadata and emits a warning", func() {
				originalDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "original",
						},
					},
				}

				currentDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "current",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalDpkg},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentDpkg},
				})

				Expect(warnings).To(ConsistOf(metadata.Warning(metadata.DebianPackageListSourceType)))

				dpkg, ok := test_utils.SelectDpkgDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(dpkg).To(Equal(currentDpkg))
			})
		})

		Context("when there is a original and there is no current", func() {
			It("retains only the empty dpkg list dependencies from the current metadata and emits a warning", func() {
				originalDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "original",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalDpkg},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{},
				})

				Expect(warnings).To(ConsistOf(metadata.Warning(metadata.DebianPackageListSourceType)))

				_, ok := test_utils.SelectDpkgDependency(result.Dependencies)
				Expect(ok).To(BeFalse())
			})
		})

		Context("when there is no original and only current", func() {
			It("retains only the dpkg list dependencies from the current metadata", func() {
				currentDpkg := metadata.Dependency{
					Type: metadata.DebianPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "current",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentDpkg},
				})

				Expect(warnings).To(BeEmpty())

				dpkg, ok := test_utils.SelectDpkgDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(dpkg).To(Equal(currentDpkg))
			})
		})
	})

	Context("rpm list dependencies on both original and current", func() {
		Context("when original and current match", func() {
			It("retains only the rpm list dependencies from the current metadata", func() {
				originalRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "some-sha256",
						},
					},
				}

				currentRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "some-sha256",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalRpm},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentRpm},
				})

				Expect(warnings).To(BeEmpty())

				rpm, ok := test_utils.SelectRpmDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(rpm).To(Equal(currentRpm))
			})
		})

		Context("when original and current don't match", func() {
			It("retains only the rpm list dependencies from the current metadata and emits a warning", func() {
				originalRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "original",
						},
					},
				}

				currentRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "current",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalRpm},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentRpm},
				})

				Expect(warnings).To(ConsistOf(metadata.Warning(metadata.RPMPackageListSourceType)))

				rpm, ok := test_utils.SelectRpmDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(rpm).To(Equal(currentRpm))
			})
		})

		Context("when there is a original and there is no current", func() {
			It("retains only the empty rpm list dependencies from the current metadata and emits a warning", func() {
				originalRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "original",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{originalRpm},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{},
				})

				Expect(warnings).To(ConsistOf(metadata.Warning(metadata.RPMPackageListSourceType)))

				_, ok := test_utils.SelectRpmDependency(result.Dependencies)
				Expect(ok).To(BeFalse())
			})
		})

		Context("when there is no original and only current", func() {
			It("retains only the rpm list dependencies from the current metadata", func() {
				currentRpm := metadata.Dependency{
					Type: metadata.RPMPackageListSourceType,
					Source: metadata.Source{
						Type: "inline",
						Version: map[string]interface{}{
							"sha256": "current",
						},
					},
				}

				result, warnings := metadata.Merge(metadata.Metadata{
					Dependencies: []metadata.Dependency{},
				}, metadata.Metadata{
					Dependencies: []metadata.Dependency{currentRpm},
				})

				Expect(warnings).To(BeEmpty())

				rpm, ok := test_utils.SelectRpmDependency(result.Dependencies)
				Expect(ok).To(BeTrue())
				Expect(rpm).To(Equal(currentRpm))
			})
		})
	})

	Context("archive dependencies on original", func() {
		It("retains the archives dependencies from the original metadata", func() {
			originalArchive := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.ArchiveType,
				},
			}
			original := metadata.Metadata{
				Dependencies: []metadata.Dependency{
					originalArchive,
					originalArchive,
				},
			}

			current := metadata.Metadata{
				Dependencies: []metadata.Dependency{},
			}

			result, warnings := metadata.Merge(original, current)
			Expect(result.Dependencies).To(Equal([]metadata.Dependency{originalArchive, originalArchive}))
			Expect(warnings).To(BeEmpty())
		})
	})

	Context("all possible types of metadata in both original and current", func() {
		It("should merge", func() {
			originalGit1 := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.GitSourceType,
					Version: map[string]interface{}{
						"commit": "commit1",
					},
				},
			}
			originalGit2 := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.GitSourceType,
					Version: map[string]interface{}{
						"commit": "commit2",
					},
				},
			}
			originalArchive2 := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.ArchiveType,
					Version: map[string]interface{}{
						"sha": "2",
					},
				},
			}
			originalArchive1 := metadata.Dependency{
				Type: metadata.PackageType,
				Source: metadata.Source{
					Type: metadata.ArchiveType,
					Version: map[string]interface{}{
						"sha": "1",
					},
				},
			}

			originalProvenance := metadata.Provenance{
				Name: "original",
			}

			original := metadata.Metadata{
				Provenance: []metadata.Provenance{originalProvenance},
				Base: metadata.Base{
					"name": "base",
				},
				Dependencies: []metadata.Dependency{
					originalArchive1,
					originalArchive2,
					{
						Type: metadata.DebianPackageListSourceType,
						Source: metadata.Source{
							Type: "inline",
							Version: map[string]interface{}{
								"sha256": "original-sha256",
							},
						},
					},
					{
						Type: metadata.RPMPackageListSourceType,
						Source: metadata.Source{
							Type: "inline",
							Version: map[string]interface{}{
								"sha256": "original-sha256",
							},
						},
					},
					originalGit1,
					originalGit2,
				},
			}

			currentBase := metadata.Base{"name": "current"}

			currentProvenance := metadata.Provenance{
				Name: "current",
			}

			currentDpkg := metadata.Dependency{
				Type: metadata.DebianPackageListSourceType,
				Source: metadata.Source{
					Type: "inline",
					Version: map[string]interface{}{
						"sha256": "current-sha256",
					},
				},
			}
			currentRpm := metadata.Dependency{
				Type: metadata.RPMPackageListSourceType,
				Source: metadata.Source{
					Type: "inline",
					Version: map[string]interface{}{
						"sha256": "current-sha256",
					},
				},
			}

			current := metadata.Metadata{
				Provenance: []metadata.Provenance{currentProvenance},
				Base:       currentBase,
				Dependencies: []metadata.Dependency{
					currentDpkg,
					currentRpm,
				},
			}

			result, warnings := metadata.Merge(original, current)

			Expect(warnings).To(ConsistOf(
				metadata.Warning(metadata.DebianPackageListSourceType),
				metadata.Warning(metadata.RPMPackageListSourceType),
				metadata.Warning("base"),
			))

			Expect(result.Provenance).To(ConsistOf(originalProvenance, currentProvenance))
			Expect(result.Base).To(Equal(currentBase))

			Expect(result.Dependencies).To(ConsistOf(
				currentDpkg,
				currentRpm,
				originalGit1,
				originalGit2,
				originalArchive1,
				originalArchive2,
			))
		})
	})
})
