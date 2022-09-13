package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	log "github.com/inconshreveable/log15"
)

var errNoGITInfoFound = fmt.Errorf("no info found (are you in a git repo ?)")

func generateBuildInfo(config *Config) (*BuildInfo, error) {
	var err error

	// The buildInfo struct is shared throughout the program to fill
	// all possible info we can get from the current build
	// and then export it to a json or template file.
	buildInfo := &BuildInfo{
		CIInfoVersion: BuildVersion,
	}

	// We get the CI info from the current CI environment
	if err = fetchCISolutionInfo(config.Directory, buildInfo); err != nil {
		return nil, fmt.Errorf("failed to fetch CI info: %w", err)
	}

	// If GIT isn't disabled, we fetch the missing information using git commands
	if err = fetchGitInfo(buildInfo, config.GitCmdMode); err != nil {
		return nil, fmt.Errorf("failed to fetch git info: %w", err)
	}

	// We fill the buildInfo struct with some information built from other parts of the struct
	if err := buildInfo.complete(); err != nil {
		return nil, fmt.Errorf("failed to complete build info: %w", err)
	}

	// At the very end we generate the version info
	if err := buildInfo.loadVersion(config); err != nil {
		return nil, fmt.Errorf("failed to load version info: %w", err)
	}

	return buildInfo, nil
}

func saveOutputFiles(config *Config, buildInfo *BuildInfo) error {
	// If requested, we export the build info to a json file
	if config.BuildInfoFile != "" {
		if err := buildInfo.save(path.Join(config.Directory, config.BuildInfoFile)); err != nil {
			return fmt.Errorf("failed to save build info: %w", err)
		}
	}

	// If request, we generate the build info file from a template
	for _, template := range config.Templates {
		if template.OutputFile != "" {
			var templateString string

			if template.InputContent != "" {
				templateString = template.InputContent
			} else {
				content, _, err := loadPathAsContent(template.InputFile, config.Directory)
				if err != nil {
					return fmt.Errorf("failed to load template file from %s: %w", template.InputFile, err)
				}
				templateString = string(content)
			}

			if err := applyTemplate(templateString, path.Join(config.Directory, template.OutputFile), buildInfo); err != nil {
				return fmt.Errorf("failed to apply template: %w", err)
			}
		}
	}

	return nil
}

func main() {
	if err := runMain(os.Args[1:]); err != nil {
		log.Error("Failed to run main", "err", err)
		os.Exit(1)
	}
}

func runMain(args []string) error {
	var config *Config

	params, err := getParams(args)

	if err != nil {
		return fmt.Errorf("failed to get params: %w", err)
	}

	// If requested, we report our version (generated by the tool itself)
	if params.Version {
		fmt.Println("Version:", BuildVersion)
		fmt.Println("Build time:", BuildDate)
		fmt.Println("Commit:", Commit)

		return nil
	}

	// If requested, we create a default config file
	if params.Init {
		if errSave := saveDefaultConfig(params); errSave != nil {
			return fmt.Errorf("failed to save default config: %w", errSave)
		}

		return nil
	}

	if params.ConfigFile == "" {
		if _, errStat := os.Stat(defaultConfigFile); errStat == nil {
			params.ConfigFile = defaultConfigFile
		}
	}

	// We load the config file
	if params.ConfigFile != "" {
		if config, err = loadConfig(params.ConfigFile); err != nil {
			return fmt.Errorf("failed to load config from \"%s\": %w", params.ConfigFile, err)
		}
	} else {
		if config, err = getEmptyConfig(); err != nil {
			return fmt.Errorf("failed to get empty config: %w", err)
		}
	}

	if params.OutputVersionFile != "" {
		config.Templates = append(config.Templates, &ConfigTemplate{
			InputContent: "{{ .Version }}",
			OutputFile:   params.OutputVersionFile,
		})
	}

	log.Debug("Loaded config", "config", config)

	// If specified, we create the build info file
	if params.OutputBuildInfoFile != "" {
		config.BuildInfoFile = params.OutputBuildInfoFile
	}

	var buildInfo *BuildInfo

	if buildInfo, err = generateBuildInfo(config); err != nil {
		log.Crit("Failed to generate build info", "err", err)

		return err
	}

	if buildInfo.CommitHash == "" {
		return errNoGITInfoFound
	}

	log.Info("Fetched build info", "buildInfo", buildInfo)

	// And then we generate all the output files
	if err = saveOutputFiles(config, buildInfo); err != nil {
		return fmt.Errorf("failed to save output files: %w", err)
	}

	return nil
}

func applyTemplate(templateString string, outputFile string, buildInfo *BuildInfo) error {
	var buffer bytes.Buffer

	if tpl, err := template.New("").Parse(templateString); err == nil {
		if errExec := tpl.Execute(&buffer, buildInfo); errExec != nil {
			return fmt.Errorf("could not execute template: %w", err)
		}
	} else {
		return fmt.Errorf("could not parse template: %w", err)
	}

	if err := os.WriteFile(outputFile, buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf("could not write output file: %w", err)
	}

	return nil
}
