package outputs

import (
	"fmt"
	"github.com/pivotal/deplab/metadata"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func WriteDpkgFile(md metadata.Metadata, dpkgFilePath string) {
	f, err := os.OpenFile(dpkgFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("no such file: %s\n", dpkgFilePath)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not close file: %s; error: %s\n", dpkgFilePath, err)
		}
	}()

	dep, err := findDpkgListInMetadata(md)
	if err != nil {
		log.Fatalf("%s", err)
	}
	pkgs := dep.Source.Metadata.(metadata.DebianPackageListSourceMetadata).Packages

	tHeader := []string{"||", "Name", "Version", "Architecture", "Description"}
	tRows := make([][]string, 0)
	tMaxLen := []int{3, 4, 7, 12, 11}

	for _, pkg := range pkgs {
		tRow := []string{"ii", pkg.Package, pkg.Version, pkg.Architecture, "Description intentionally left blank"}
		for i, v := range tMaxLen {
			if utf8.RuneCountInString(tRow[i]) > v {
				tMaxLen[i] = utf8.RuneCountInString(tRow[i])
			}
		}
		tRows = append(tRows, tRow)
	}

	header := `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
`

	pRows := header
	unpaddedFmtString := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds\n", tMaxLen[0]+1, tMaxLen[1]+1, tMaxLen[2]+1, tMaxLen[3]+1, tMaxLen[4]+1)
	fmtString := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n", tMaxLen[0], tMaxLen[1], tMaxLen[2], tMaxLen[3], tMaxLen[4])

	pRows += fmt.Sprintf(fmtString, tHeader[0], tHeader[1], tHeader[2], tHeader[3], tHeader[4])
	pRows += fmt.Sprintf(unpaddedFmtString,
		"+++-",
		fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[1])),
		fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[2])),
		fmt.Sprintf("%s-", strings.Repeat("=", tMaxLen[3])),
		fmt.Sprintf("%s", strings.Repeat("=", tMaxLen[4])),
	)

	for _, v := range tRows {
		pRows += fmt.Sprintf(fmtString, v[0], v[1], v[2], v[3], v[4])
	}

	contents := []byte(pRows)
	_, err = f.Write(contents)
	if err != nil {
		log.Fatalf("Could not write to file: %s\n", dpkgFilePath)
	}
}

func findDpkgListInMetadata(md metadata.Metadata) (metadata.Dependency, error) {
	for _, v := range md.Dependencies {
		if v.Type == "debian_package_list" {
			return v, nil
		}
	}

	return metadata.Dependency{}, fmt.Errorf("could not find debian_package_list in metadata")
}
