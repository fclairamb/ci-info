package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputVersionFile(t *testing.T) {
	a := require.New(t)
	version, err := getVersionFromFile("testdata/README.md", "Version: [0-9+\\.]+")
	a.Nil(err)
	a.Equal("1.2.3", version)
}

func TestInputVersion(t *testing.T) {
	a := require.New(t)
	version, err := getVersionFromContent(`
This is about
Version: 1.2.3

`, "Version: ([0-9+\\.]+)")
	a.Nil(err)
	a.Equal("1.2.3", version)
}
