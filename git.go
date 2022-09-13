package main

import (
	"os"
	"path"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// var repo *git.Repository

func getRepo(dir string) (*git.Repository, error) {
	if dir == "" {
		var errCwd error
		dir, errCwd = os.Getwd()

		if errCwd != nil {
			return nil, errCwd
		}
	}

	for i := 0; i < 20; i++ {
		if st, errSt := os.Stat(dir + "/.git"); errSt == nil {
			if st.IsDir() {
				break
			}
		} else if !os.IsNotExist(errSt) {
			return nil, errSt
		}

		dir = path.Dir(dir)
	}

	return git.PlainOpen(dir)
}

func getHead() (*git.Repository, *plumbing.Reference, error) {
	repo, err := getRepo("")
	if err != nil {
		return nil, nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, nil, err
	}

	return repo, ref, nil
}

func fetchGitInfo(info *BuildInfo) error {
	// If we have everything, we don't need to fetch anything
	if info.GitCommitHash != "" && info.GitCommitDate != "" && (info.GitBranch != "" || info.GitTag != "") {
		return nil
	}

	repo, ref, err := getHead()
	if err != nil {
		return err
	}

	if info.GitCommitHash == "" {
		info.GitCommitHash = ref.Hash().String()
	}

	if info.GitBranch == "" {
		info.GitBranch = ref.Name().Short()
	}

	if info.GitCommitDate == "" {
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		info.GitCommitDate = commit.Committer.When.Format("2006-01-02 15:04:05 -0700")
	}

	if info.GitTag == "" {
		tags, err := repo.Tags()
		if err != nil {
			return err
		}

		err = tags.ForEach(func(t *plumbing.Reference) error {
			if t.Hash() == ref.Hash() {
				info.GitTag = t.Name().Short()
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if info.GitLastTag == "" {
		lastTag, err := getLastTagFromRepository(repo)

		if err != nil {
			return err
		}

		info.GitLastTag = lastTag
	}

	return nil
}

// See: https://github.com/src-d/go-git/issues/1030
func getLastTagFromRepository(repository *git.Repository) (string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, subErr := repository.ResolveRevision(revision)
		if subErr != nil {
			return subErr
		}

		commit, subErr := repository.CommitObject(*tagCommitHash)
		if subErr != nil {
			return subErr
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return latestTagName, nil
}
