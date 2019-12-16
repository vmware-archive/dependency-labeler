package rpm

import (
	"reflect"
	"strings"

	"github.com/pivotal/deplab/pkg/metadata"
)

const sep = "\t"

// UnmarshalPackage matches the line items to the index of the struct field
func UnmarshalPackage(packageLine string) metadata.RpmPackage {
	rpmPackage := metadata.RpmPackage{}
	rpmPackageValue := reflect.ValueOf(&rpmPackage).Elem()

	values := strings.Split(packageLine, sep)
	for i, value := range values {
		rpmPackageValue.Field(i).SetString(value)
	}
	return rpmPackage
}

// QueryFormat uses the struct definition to create a queryformat for rpm
func QueryFormat() string {
	var fields []string

	rpmPackage := metadata.RpmPackage{}
	rpmPackageValue := reflect.ValueOf(&rpmPackage).Elem()
	rpmPackageType := rpmPackageValue.Type()

	for i := 0; i < rpmPackageValue.NumField(); i++ {
		field := rpmPackageType.Field(i).Tag.Get("rpm")
		fields = append(fields, "%{"+field+"}")
	}

	return strings.Join(fields, sep) + "\n"
}
