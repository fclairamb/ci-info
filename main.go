package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/inconshreveable/log15"
)

// CLI parameters
type CmdParams struct {
	ConfigFile    string
	Version       string
	BuildInfoFile string
	Init          bool
}

// ConfigVersionInputFile defines how we shall fetch the input version
type ConfigVersionInputFile struct {
	File    string `json:"file"`
	Pattern string `json:"pattern"`
}

type ConfigVersionInputTag struct {
	Pattern string `json:"pattern"`
}

// ConfigTemplate
type ConfigTemplate struct {
	InputFile  string `json:"input_file"`
	OutputFile string `json:"output_file"`
}

// Config defines the configuration for ci-info
type Config struct {
	InputVersionFile ConfigVersionInputFile `json:"version_input_file"`
	InputVersionTag  ConfigVersionInputTag  `json:"version_input_tag"`
	Template         ConfigTemplate         `json:"template"`
	BuildInfoFile    string                 `json:"build_info_file"`
	DisableGitCmd    bool                   `json:"disabled_git_cmd"`
}

type gitInfoFetch struct {
	info    *string
	command []string
}

func loadConfig(fileName string) (*Config, error) {
	jsonContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(jsonContent, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadVersion(fileName string, pattern string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	version := re.FindString(string(content))
	return version, nil
}

func createDefaultConfig() *Config {
	return &Config{
		InputVersionFile: ConfigVersionInputFile{
			File:    "README.md",
			Pattern: "Version: ([0-9.]+)\n",
		},
		InputVersionTag: ConfigVersionInputTag{
			Pattern: "^v?([0-9.]+)$",
		},
		Template: ConfigTemplate{
			InputFile:  "build.go.tpl",
			OutputFile: "build.go",
		},
		BuildInfoFile: "build.json",
	}
}

func saveDefaultConfig(params *CmdParams) error {
	config := createDefaultConfig()
	jsonContent, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(params.ConfigFile, jsonContent, 0644)
}

func getVersionFromTag(tag string, config *Config) (string, error) {
	if tag == "" || config.InputVersionTag.Pattern == "" {
		return "", nil
	}
	re, err := regexp.Compile(config.InputVersionTag.Pattern)
	if err != nil {
		return "", fmt.Errorf("could not compile pattern: %w", err)
	}

	return re.FindString(tag), nil
}

func getVersionFromFile(config *Config) (string, error) {
	if config.InputVersionFile.File == "" || config.InputVersionFile.Pattern == "" {
		return "", nil
	}
	content, err := ioutil.ReadFile(config.InputVersionFile.File)
	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
	re, err := regexp.Compile(config.InputVersionFile.Pattern)
	if err != nil {
		return "", fmt.Errorf("could not compile pattern: %w", err)
	}

	return re.FindString(string(content)), nil
}

func main() {
	params := &CmdParams{}
	flag.StringVar(&params.ConfigFile, "c", ".ci-info.json", "config file")
	flag.StringVar(&params.Version, "v", "", "version")
	flag.StringVar(&params.BuildInfoFile, "b", "build.json", "build info file")
	flag.BoolVar(&params.Init, "init", false, "init config file")
	flag.Parse()

	var config *Config
	var err error

	if params.Init {
		err = saveDefaultConfig(params)
		if err != nil {
			log.Crit("could not save default config", "err", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	if config, err = loadConfig(params.ConfigFile); err != nil {
		log.Crit("Failed to load config from file", "file", params.ConfigFile, "err", err)
		os.Exit(1)
	}

	log.Debug("Loaded config", "config", config)

	buildInfo := &BuildInfo{}

	if params.Version != "" {
		buildInfo.Version = params.Version
	}

	if params.BuildInfoFile != "" {
		config.BuildInfoFile = params.BuildInfoFile
	}

	fetchers := []CIInfoFetcher{
		&CircleCIInfoFetcher{},
		&GithubActionsCIInfoFetcher{},
		&GitLabInfoFetcher{},
		&DroneCIInfoFetcher{},
		&TravisCIInfoFetcher{},
		&JenkinsCIInfoFetcher{},
	}
	for _, fetcher := range fetchers {
		if !fetcher.Detect() {
			continue
		}
		log.Info("Found CI info fetcher", "fetcher", fetcher.String())
		if err = fetcher.Fetch(buildInfo); err != nil {
			log.Error("Failed to fetch CI info", "fetcher", fetcher, "err", err)
			os.Exit(1)
		}
		break
	}

	if !config.DisableGitCmd {
		if err = fetchGitInfoWithCmd(buildInfo); err != nil {
			log.Crit("Failed to fetch git info", "err", err)
			os.Exit(1)
		}
	}

	if buildInfo.CommitHash == "" {
		log.Error("Could not find anything to use")
		os.Exit(1)
	}

	var tagVersion, fileVersion string

	if tagVersion, err = getVersionFromTag(buildInfo.CommitTag, config); err != nil {
		log.Error("Failed to get version from tag", "err", err)
		os.Exit(1)
	}

	if fileVersion, err = getVersionFromFile(config); err != nil {
		log.Error("Failed to get version from file", "err", err)
		os.Exit(1)
	}

	if tagVersion != "" {
		buildInfo.Version = tagVersion
	} else if fileVersion != "" {
		buildInfo.Version = fileVersion
	}

}
