package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteConfig(t *testing.T) {
	a := assert.New(t)

	params := &CmdParams{
		ConfigFile: "https://raw.githubusercontent.com/fclairamb/ci-info/main/testdata/c/.ci-info.json",
	}

	conf, err := loadConfig(params.ConfigFile)
	a.NoError(err)
	a.NotNil(conf)

	a.Equal("VERSION", conf.InputVersionFile.File)
	a.Equal("build.json", conf.BuildInfoFile)
}
