package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"

	log "github.com/inconshreveable/log15"
)

// ConfigVersionInputFile defines how we shall fetch the input version
type ConfigVersionInputFile struct {
	File    string `json:"file"`
	Pattern string `json:"pattern"`
}

// ConfigVersionInputTag defines how we shall fetch the input version
type ConfigVersionInputTag struct {
	Pattern string `json:"pattern"`
}

// ConfigTemplate defines the template configuration
type ConfigTemplate struct {
	InputFile    string `json:"input_file"`
	InputContent string `json:"input_content"`
	OutputFile   string `json:"output_file"`
}

// Config defines the configuration for ci-info
type Config struct {
	InputVersionFile ConfigVersionInputFile `json:"version_input_file"`
	InputVersionTag  ConfigVersionInputTag  `json:"version_input_tag"`
	Templates        []*ConfigTemplate      `json:"templates"`
	BuildInfoFile    string                 `json:"build_info_file"`
	GitCmdMode       bool                   `json:"git_cmd_mode"`
}

const defaultConfigFile = ".ci-info.json"

var regexURL = regexp.MustCompile(`^https?://`)

func loadPathAsContent(path string) ([]byte, error) {
	var reader io.ReadCloser

	if regexURL.MatchString(path) {
		resp, err := http.Get(path) //nolint:gosec
		if err != nil {
			return nil, err
		}

		reader = resp.Body

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Crit("could not close response body", "err", err)
			}
		}()
	} else {
		file, err := os.Open(path) //nolint:gosec
		if err != nil {
			return nil, err
		}

		reader = file
	}

	return io.ReadAll(reader)
}

func loadConfig(fileName string) (*Config, error) {
	jsonContent, err := loadPathAsContent(fileName)
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

func createDefaultConfig() *Config {
	return &Config{
		InputVersionFile: ConfigVersionInputFile{
			File:    "README.md",
			Pattern: "Version: ([0-9.]+)\n",
		},
		InputVersionTag: ConfigVersionInputTag{
			Pattern: "^v?([0-9.]+)$",
		},
		Templates: []*ConfigTemplate{{
			InputFile:  "build.go.tpl",
			OutputFile: "build.go",
		}},
		BuildInfoFile: "build.json",
	}
}

func saveDefaultConfig(params *CmdParams) error {
	config := createDefaultConfig()
	jsonContent, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(params.ConfigFile, jsonContent, 0600)
}
