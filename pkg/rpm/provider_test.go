package rpm_test

import (
	"os"

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

func (m MockImage) AbsolutePath(string) string {
	path, err := filepath.Abs(m.path)

	Expect(err).ToNot(HaveOccurred())
	return path
}

var _ = Describe("Pkg/Rpm/Provider", func() {
	It("should generate list of dependencies", func() {
		md, err := rpm.BuildRPMMetadata(MockImage{"../../test/integration/assets/rpm"})

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
			}))
		}
	})

	It("returns an error if no rpm package is found", func() {
		_, err := rpm.BuildRPMMetadata(MockImage{"temp/this-path-does-not-exists"})

		Expect(err).To(MatchError(ContainSubstring("no rpm packages data found")))
	})

	It("returns an error if rpm is not in the PATH", func() {
		PATH := os.Getenv("PATH")
		Expect(os.Setenv("PATH", "")).ToNot(HaveOccurred())

		defer func() {
			Expect(os.Setenv("PATH", PATH)).ToNot(HaveOccurred())
		}()

		_, err := rpm.BuildRPMMetadata(MockImage{"../../test/integration/assets/rpm"})

		Expect(err).To(MatchError(
			SatisfyAll(
				ContainSubstring("executable file not found in $PATH"),
				ContainSubstring("rpm"))))

	})
})

/*
 md, err := rpm.BuildRPMMetadata(MockImage{"this-path-does-not-exists"})
	if err != nil {
		log.warn(err)
	}else{
		append()
	}


	deplab -i photon --experimental-rpm --output-tar /tmp/image.tar
	error: you should have rpm in PATH
	warn:  you should have rpm in PATH


	deplab -i ubuntu --experimental-rpm --output-tar /tmp/image.tar
	warn: this images is not rpm distro


	deplab -i alpine  --output-tar /tmp/image.tar
	warn: this images does not have dpkg db
	warn: this images does not have rpm db
	warn: this images does not have pacman db

if ./rpm folder === false and no rpm cli then
	ignore
else if ./rpm folder == true and no rpm cli then
	error!!!


	if packages, err := getDPKG; err == nil {
		append(packages, )
	}
	if packages, err := getRPM; err == nil {
		append(packages, )
	}


*/
