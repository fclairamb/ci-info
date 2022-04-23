package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDate(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		CommitDate: "2022-04-21 11:52:09 +0200",
	}
	a.Nil(bi.complete())
	a.Equal("2022-04-21-0952", bi.CommitDateClean)
}

func TestTagOrBranch(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		CommitBranch: "feature/cool-one",
		CommitHash:   "a6850c90c8d3c81377cee5701f79dfbbd6e5a756",
	}
	a.Nil(bi.complete())
	a.Equal("feature-cool-one", bi.CommitRef)
	a.Equal("feature-cool-one-a6850c9", bi.CommitSmart)

	bi.CommitTag = "v1.2.3"
	a.Nil(bi.complete())
	a.Equal("v1.2.3", bi.CommitRef)
	a.Equal("v1.2.3", bi.CommitSmart)
}

func TestHash(t *testing.T) {
	a := require.New(t)
	bi := &BuildInfo{
		CommitHash: "a6850c90c8d3c81377cee5701f79dfbbd6e5a756",
	}
	a.Nil(bi.complete())
	a.Equal("a6850c9", bi.CommitHashShort)
	a.Equal("a6850c9", bi.CommitSmart)
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
		CommitHash:      "abcdef",
		CommitDate:      "2022-04-21 11:52:09 +0200",
		CommitBranch:    "main",
	}
	a.Nil(bi.complete())
	a.Nil(bi.save("/tmp/buildinfo.json"))
	a.NotNil(bi.save("/not-existing-path/buildinfo.json"))
}
