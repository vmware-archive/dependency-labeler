package metadata

import (
	"reflect"
)

type Warning string

func Merge(original, current Metadata) (Metadata, []Warning) {
	var warnings []Warning
	var newDependencies []Dependency

	if len(original.Base) > 0 && !reflect.DeepEqual(original.Base, current.Base) {
		warnings = append(warnings, "base")
	}

	newDependencies, warnings = selectAdditionalDependencies(DebianPackageListSourceType, newDependencies, warnings, original, current)
	newDependencies, warnings = selectAdditionalDependencies(RPMPackageListSourceType, newDependencies, warnings, original, current)
	newDependencies, warnings = selectAdditionalDependencies(BuildpackMetadataType, newDependencies, warnings, original, current)

	for _, dep := range original.Dependencies {
		if dep.Source.Type == GitSourceType {
			newDependencies = append(newDependencies, dep)
		} else if dep.Source.Type == ArchiveType {
			newDependencies = append(newDependencies, dep)
		}
	}

	return Metadata{
		Provenance:   append(original.Provenance, current.Provenance...),
		Base:         current.Base,
		Dependencies: newDependencies,
	}, warnings
}

func selectAdditionalDependencies(sourceType string, dependencies []Dependency, warnings []Warning, original Metadata, current Metadata) ([]Dependency, []Warning) {
	originalDpkgDependency, dpkgPresentInSource := SelectDependency(original.Dependencies, sourceType)
	currentDpkgDependency, dpkgPresentInAdditional := SelectDependency(current.Dependencies, sourceType)

	if dpkgPresentInSource && !reflect.DeepEqual(originalDpkgDependency, currentDpkgDependency) {
		warnings = append(warnings, Warning(sourceType))
	}

	if dpkgPresentInAdditional {
		dependencies = append(dependencies, currentDpkgDependency)
	}

	return dependencies, warnings
}

func SelectDependency(dependencies []Dependency, dependencyType string) (Dependency, bool) {
	for _, dependency := range dependencies {
		if dependency.Type == dependencyType {
			return dependency, true
		}
	}
	return Dependency{}, false
}
