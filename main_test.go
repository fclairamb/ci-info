package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteConfig(t *testing.T) {
	a := assert.New(t)

	params, err := getParams([]string{
		"-c", "https://raw.githubusercontent.com/fclairamb/ci-info/main/testdata/c/.ci-info.json",
	})
	a.NoError(err)
	a.NotNil(params)

	conf, err := loadConfig(params.ConfigFile)
	a.NoError(err)
	a.NotNil(conf)

	a.Equal("VERSION", conf.InputVersionFile.File)
	a.Equal("build.json", conf.BuildInfoFile)
}

func TestRunOutputVersionFile(t *testing.T) {
	a := assert.New(t)

	outputFile := "testdata/output-version-file.txt.out"

	_ = os.Remove(outputFile)
	a.NoFileExists(outputFile)

	a.NoError(runMain([]string{
		"-vf",
		outputFile,
	}))

	a.FileExists(outputFile)
}

func TestRunOutputBuildInfoFile(t *testing.T) {
	a := assert.New(t)

	outputFIle := "testdata/build.json.out"

	_ = os.Remove(outputFIle)
	a.NoFileExists(outputFIle)

	a.NoError(runMain([]string{
		"-b",
		outputFIle,
	}))

	a.FileExists(outputFIle)
}

func TestRunStandardOne(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{}))
}
