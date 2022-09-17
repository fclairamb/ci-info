package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubActions(t *testing.T) {
	a := assert.New(t)

	config := createDefaultConfig()

	t.Run("tag", func(t *testing.T) {
		bi := &BuildInfo{}
		t.Setenv("GITHUB_ACTION", "true")
		t.Setenv("GITHUB_REF", "refs/tags/v1.2.3")

		a.NoError(fetchCISolutionInfo("", bi))

		a.Equal("github-actions", bi.CISolution)
		a.Equal("v1.2.3", bi.GitTag)

		a.NoError(bi.complete())
		a.Equal("v1.2.3", bi.GitRef)

		a.NoError(bi.loadVersion(config))
		a.Equal("1.2.3", bi.Version)
	})

	t.Run("branch", func(t *testing.T) {
		bi := &BuildInfo{}
		t.Setenv("GITHUB_ACTION", "true")
		t.Setenv("GITHUB_REF", "refs/heads/feature/branch")

		a.NoError(fetchCISolutionInfo("", bi))

		a.Equal("github-actions", bi.CISolution)
		a.Equal("feature/branch", bi.GitBranch)

		a.NoError(bi.complete())

		a.Equal("", bi.GitTag)
	})
}
