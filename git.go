package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/inconshreveable/log15"
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

func fetchGitInfoNative(info *BuildInfo) error {
	// If we have everything, we don't need to fetch anything
	if info.CommitHash != "" && info.CommitDate != "" && (info.CommitBranch != "" || info.CommitTag != "") {
		return nil
	}

	repo, ref, err := getHead()
	if err != nil {
		return err
	}

	if info.CommitHash == "" {
		info.CommitHash = ref.Hash().String()
	}

	if info.CommitBranch == "" {
		info.CommitBranch = ref.Name().Short()
	}

	if info.CommitDate == "" {
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		info.CommitDate = commit.Committer.When.Format("2006-01-02 15:04:05 -0700")
	}

	if info.CommitTag == "" {
		tags, err := repo.Tags()
		if err != nil {
			return err
		}

		err = tags.ForEach(func(t *plumbing.Reference) error {
			if t.Hash() == ref.Hash() {
				info.CommitTag = t.Name().Short()
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchGitInfoWithCmd(info *BuildInfo) error {
	type gitInfoFetch struct {
		info    *string
		command []string
	}

	var gitCommands = []gitInfoFetch{
		{&info.CommitTag, []string{"git", "tag", "--points-at", "HEAD"}},
		{&info.CommitHash, []string{"git", "rev-parse", "HEAD"}},
		{&info.CommitBranch, []string{"git", "branch", "--show-current"}},
		{&info.CommitDate, []string{"git", "show", "--quiet", "--format=%ci"}},
	}

	for i := range gitCommands {
		gifc := &gitCommands[i]
		if *gifc.info != "" {
			continue
		}

		log.Debug("Fetching git info", "command", gifc.command)
		cmd := exec.Command(gifc.command[0], gifc.command[1:]...) //nolint:gosec

		var outBin []byte
		var err error

		if outBin, err = cmd.CombinedOutput(); err != nil {
			return err
		}

		out := strings.TrimRight(string(outBin), "\n")
		log.Debug("Fetched git info", "command", gifc.command, "output", out)
		*gifc.info = out
	}

	return nil
}
