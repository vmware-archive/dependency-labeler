package git

import (
	"fmt"

	"github.com/pivotal/deplab/pkg/metadata"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const SourceType = "git"

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
			Type: SourceType,
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
