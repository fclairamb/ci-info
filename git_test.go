package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchGitInfoWithCmd(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{}
	a.NoError(fetchGitInfoWithCmd(bi))
	testGitInfo(a, bi)
}

func TestFetchGitInfoNative(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{}
	a.NoError(fetchGitInfoNative(bi))
	testGitInfo(a, bi)
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
