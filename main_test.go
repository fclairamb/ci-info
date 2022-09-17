package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteConfig(t *testing.T) {
	a := assert.New(t)

	params, err := getParams([]string{
		"-c", "https://raw.githubusercontent.com/fclairamb/ci-info/main/samples/c/.ci-info.json",
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

func TestMainRunMaven(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-c", "testdata/maven/.ci-info.json"}))
	a.FileExists("testdata/maven/build.json")
	a.FileExists("testdata/maven/version.txt")
}

func TestMainRunNpm(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-c", "testdata/npm/.ci-info.json"}))
	a.FileExists("testdata/npm/build.json")
	a.FileExists("testdata/npm/version.txt")
}

func TestRunStandardOne(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{}))
}

func TestVersionFromLastTag(t *testing.T) {
	a := assert.New(t)

	config, err := getEmptyConfig()
	a.NoError(err)
	bi, err := generateBuildInfo(config)
	a.NoError(err)
	a.NotEmpty(bi.GitLastTag)
	a.NotEmpty(bi.Version)
}

func TestMainVersion(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-v"}))
}

func TestMainHelp(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-h"}))
	a.NoError(runMain([]string{"-help"}))
	a.NoError(runMain([]string{"--help"}))
	a.Error(runMain([]string{"-unknown-flag"}))
}

func TestMainInit(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-i", "-c", "testdata/.ci-info.init.out"}))
	a.FileExists("testdata/.ci-info.init.out")
}

func TestMainLoggingLevel(t *testing.T) {
	a := assert.New(t)

	a.NoError(runMain([]string{"-l", "debug"}))
	a.NoError(runMain([]string{"-l", "info"}))
	a.NoError(runMain([]string{"-l", "warn"}))
	a.NoError(runMain([]string{"-l", "error"}))
	a.Error(runMain([]string{"-l", "bad"}))
}
