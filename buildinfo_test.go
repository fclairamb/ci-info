package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDate(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		GitCommitDate: "2022-04-21 11:52:09 +0200",
	}
	a.Nil(bi.complete())
	a.Equal("2022-04-21-0952", bi.GitCommitDateClean)
}

func TestTagOrBranch(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		GitBranch:     "feature/cool-one",
		GitCommitHash: "a6850c90c8d3c81377cee5701f79dfbbd6e5a756",
	}
	a.Nil(bi.complete())
	a.Equal("feature-cool-one", bi.GitRef)
	a.Equal("feature-cool-one-a6850c9", bi.GitSmartRef)

	bi.GitTag = "v1.2.3"
	a.Nil(bi.complete())
	a.Equal("v1.2.3", bi.GitRef)
	a.Equal("v1.2.3", bi.GitSmartRef)
}

func TestHash(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		GitCommitHash: "a6850c90c8d3c81377cee5701f79dfbbd6e5a756",
	}
	a.Nil(bi.complete())
	a.Equal("a6850c9", bi.GitCommitHashShort)
	a.Equal("a6850c9", bi.GitSmartRef)
}

func TestVersionAuto(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		VersionDeclared: "1.2.3",
	}
	a.Nil(bi.complete())
	a.Equal("1.2.3", bi.VersionDeclared)
}

func TestSave(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		VersionDeclared: "1.2.3",
		GitCommitHash:   "a6850c90c8d3c81377cee5701f79dfbbd6e5a756",
		GitCommitDate:   "2022-04-21 11:52:09 +0200",
		GitBranch:       "main",
	}
	a.Nil(bi.complete())
	a.Nil(bi.save("/tmp/buildinfo.json"))
	a.NotNil(bi.save("/not-existing-path/buildinfo.json"))
}

func TestReadme(t *testing.T) {
	a := require.New(t)
	config := &Config{
		Templates: []*ConfigTemplate{{
			InputFile:  "README.md",
			OutputFile: "/tmp/README.md",
		}},
		InputVersionFile: ConfigVersionInputFile{
			File:    "testdata/README.md",
			Pattern: "Version: ([0-9+\\.]+)",
		},
	}
	bi, err := generateBuildInfo(config)
	a.NoError(err)
	a.NotNil(bi)
	a.NotEmpty(bi.Version)

	a.NoError(saveOutputFiles(config, bi))
}
