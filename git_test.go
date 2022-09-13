package main

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchGitInfo(t *testing.T) {
	a := require.New(t)

	bi := &BuildInfo{}
	a.NoError(fetchGitInfo(bi))
	testGitInfo(a, bi)
}

func TestGitInSubpath(t *testing.T) {
	a := require.New(t)

	dir, err := os.Getwd()
	a.NoError(err)

	dir = path.Join(dir, "testdata")

	repo, err := getRepo(dir)
	a.NoError(err)
	a.NotNil(repo)
}

func TestGitLastTag(t *testing.T) {
	a := require.New(t)

	bi := &BuildInfo{}
	a.NoError(fetchGitInfo(bi))
	a.NotEmpty(bi.GitLastTag)
	a.True(regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(bi.GitLastTag))
}

func testGitInfo(a *require.Assertions, bi *BuildInfo) {
	a.NotEmpty(bi.GitCommitHash)

	// Github creates a detached branch for PRs and this prevents from detecting a branch:
	if os.Getenv("GITHUB_ACTION") == "" {
		a.NotEmpty(bi.GitBranch)
	}

	a.NotEmpty(bi.GitCommitDate)
	a.Regexp("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} \\+[0-9]{4}", bi.GitCommitDate)

	a.Nil(bi.complete())

	if os.Getenv("GITHUB_ACTION") == "" {
		a.NotEmpty(bi.GitBranchClean)
	}

	a.NotEmpty(bi.GitCommitDateClean)
	a.Regexp("[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{4}", bi.GitCommitDateClean)
}
