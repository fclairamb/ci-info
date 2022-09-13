package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchGitInfo(t *testing.T) {
	a := require.New(t)

	t.Run("native", func(t *testing.T) {
		bi := &BuildInfo{}
		a.NoError(fetchGitInfo(bi, false))
		testGitInfo(a, bi)
	})

	t.Run("cmd", func(t *testing.T) {
		bi := &BuildInfo{}
		a.NoError(fetchGitInfo(bi, true))
		testGitInfo(a, bi)
	})
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

func testGitInfo(a *require.Assertions, bi *BuildInfo) {
	a.NotEmpty(bi.CommitHash)

	// Github creates a detached branch for PRs and this prevents from detecting a branch:
	if os.Getenv("GITHUB_ACTION") == "" {
		a.NotEmpty(bi.CommitBranch)
	}

	a.NotEmpty(bi.CommitDate)
	a.Regexp("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} \\+[0-9]{4}", bi.CommitDate)

	a.Nil(bi.complete())

	if os.Getenv("GITHUB_ACTION") == "" {
		a.NotEmpty(bi.CommitBranchClean)
	}

	a.NotEmpty(bi.CommitDateClean)
	a.Regexp("[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{4}", bi.CommitDateClean)
}
