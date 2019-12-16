package rpm_test

import (
	"os"

	"github.com/pivotal/deplab/pkg/common"

	"github.com/onsi/gomega/gstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/deplab/pkg/metadata"
	"github.com/pivotal/deplab/pkg/rpm"

	"path/filepath"
)

type MockImage struct {
	path string
}

func (m MockImage) Cleanup() {
	panic("implement me")
}

func (m MockImage) GetFileContent(string) (string, error) {
	panic("implement me")
}

func (m MockImage) GetDirContents(string) ([]string, error) {
	panic("implement me")
}

func (m MockImage) AbsolutePath(string) (string, error) {
	path, err := filepath.Abs(m.path)

	Expect(err).ToNot(HaveOccurred())
	return path, err
}

var _ = Describe("Pkg/Rpm/Provider", func() {
	var (
		rpmProvider common.Provider
	)
	BeforeEach(func() {
		rpmProvider = rpm.RPMProvider{}
	})

	//rpm leaves __db.001 etc. files in the folder when it runs; we should try to clean those up
	AfterEach(func() {
		files, err := filepath.Glob("../../test/integration/assets/rpm/__db.*")
		Expect(err).ToNot(HaveOccurred())
		for _, f := range files {
			err := os.Remove(f)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("should generate list of dependencies", func() {
		md, err := rpmProvider.BuildDependencyMetadata(MockImage{"../../test/integration/assets/rpm"})

		Expect(err).ToNot(HaveOccurred())
		packages := md.Source.Metadata.(metadata.RpmPackageListSourceMetadata).Packages
		Expect(packages).To(HaveLen(34))

		for _, p := range packages {
			Expect(p.Package).ToNot(BeEmpty())
			Expect(p).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Package":      Not(BeEmpty()),
				"Version":      Not(BeEmpty()),
				"License":      Not(BeEmpty()),
				"Architecture": Not(BeEmpty()),
				"SourceRpm":    Not(BeEmpty()),
			}))
		}
	})

	It("returns an empty struct if no rpm database is found", func() {
		tempDirPath := "/tmp/this-path-does-not-exists"
		packages, err := rpmProvider.BuildDependencyMetadata(MockImage{tempDirPath})
		Expect(err).NotTo(HaveOccurred())

		Expect(packages).To(Equal(metadata.Dependency{}))

		_ = os.Remove(tempDirPath)
	})

	It("returns an error if rpm is not in the PATH", func() {
		PATH := os.Getenv("PATH")
		Expect(os.Setenv("PATH", "")).ToNot(HaveOccurred())

		defer func() {
			Expect(os.Setenv("PATH", PATH)).ToNot(HaveOccurred())
		}()

		_, err := rpmProvider.BuildDependencyMetadata(MockImage{"../../test/integration/assets/rpm"})

		Expect(err).To(MatchError(SatisfyAll(
			ContainSubstring("an rpm database exists at"),
			ContainSubstring("but rpm is not installed and available on your path"))))

	})
})
