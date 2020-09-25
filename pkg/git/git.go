// Copyright (c) 2019-2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

package git

import (
	"fmt"

	"github.com/vmware-tanzu/dependency-labeler/pkg/common"

	"github.com/vmware-tanzu/dependency-labeler/pkg/image"

	"github.com/vmware-tanzu/dependency-labeler/pkg/metadata"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func Provider(dli image.Image, params common.RunParams, md metadata.Metadata) (metadata.Metadata, error) {
	for _, path := range params.GitPaths {
		dependency, err := BuildDependencyMetadata(path)
		if err != nil {
			return metadata.Metadata{}, err
		}
		md.Dependencies = append(md.Dependencies, dependency)
	}

	return md, nil
}

func BuildDependencyMetadata(pathToGit string) (metadata.Dependency, error) {
	repo, err := git.PlainOpen(pathToGit)
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("cannot open git repository \"%s\": %s\n", pathToGit, err)
	}

	ref, err := repo.Head()
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("cannot find head of git repository: %s\n", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("cannot find remotes for repository: %s\n", err)
	}

	var refs []string
	tags, err := repo.Tags()
	if err != nil {
		return metadata.Dependency{}, fmt.Errorf("error finding tags: %s\n", err)
	}

	tags.ForEach(func(tagRef *plumbing.Reference) error {
		if tagRef.Type() == plumbing.HashReference && tagRef.Hash() == ref.Hash() {
			refs = append(refs, tagRef.Name().Short())
		}

		return nil
	})

	return metadata.Dependency{
		Type: "package",
		Source: metadata.Source{
			Type: metadata.GitSourceType,
			Version: map[string]interface{}{
				"commit": ref.Hash().String(),
			},
			Metadata: metadata.GitSourceMetadata{
				URL:  remotes[0].Config().URLs[0],
				Refs: refs,
			},
		},
	}, nil
}
