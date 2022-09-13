package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	log "github.com/inconshreveable/log15"
)

// BuildInfo contains all the information about the build
type BuildInfo struct {
	CIInfoVersion      string `json:"ci_info_version"`
	VersionDeclared    string `json:"-"`
	Version            string `json:"version,omitempty"`
	GitCommitHash      string `json:"git_hash,omitempty"`
	GitCommitHashShort string `json:"-"`
	GitCommitDate      string `json:"git_date,omitempty"`
	GitCommitDateClean string `json:"-"`
	GitBranch          string `json:"git_branch,omitempty"`
	GitBranchClean     string `json:"-"`
	GitTag             string `json:"git_tag,omitempty"`
	GitRef             string `json:"-"`
	GitSmartRef        string `json:"-"`
	GitLastTag         string `json:"-"`
	BuildDate          string `json:"build_date,omitempty"`
	BuildHost          string `json:"build_host,omitempty"`
	BuildUser          string `json:"build_user,omitempty"`
	CISolution         string `json:"ci_solution,omitempty"`
	CIBuildNumber      string `json:"ci_build_number,omitempty"`
	PackageManager     string `json:"package_manager,omitempty"`
}

var reBranchClean = regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)

const timeFormat = "2006-01-02-1504"

func (bi *BuildInfo) complete() error {
	if bi.GitCommitHash != "" {
		bi.GitCommitHashShort = bi.GitCommitHash[:7]
	}

	if bi.GitCommitDate != "" {
		date, err := time.Parse("2006-01-02 15:04:05 -0700", bi.GitCommitDate)
		if err != nil {
			return fmt.Errorf("could not parse commit date: %w", err)
		}

		bi.GitCommitDateClean = date.UTC().Format(timeFormat)
	}

	if bi.GitBranch != "" {
		bi.GitBranchClean = reBranchClean.ReplaceAllString(bi.GitBranch, "-")
	}

	switch {
	case bi.GitTag != "":
		bi.GitRef = bi.GitTag
		bi.GitSmartRef = bi.GitTag
	case bi.GitBranchClean != "":
		bi.GitRef = bi.GitBranchClean
		bi.GitSmartRef = bi.GitBranchClean + "-" + bi.GitCommitHashShort
	default:
		bi.GitSmartRef = bi.GitCommitHashShort
	}

	if bi.BuildHost == "" {
		if host, err := os.Hostname(); err != nil {
			log.Warn("Could not get hostname", "err", err)
		} else {
			bi.BuildHost = host
		}
	}

	if bi.BuildUser == "" {
		bi.BuildUser = os.Getenv("USER")
	}

	if bi.BuildDate == "" {
		bi.BuildDate = time.Now().UTC().Format(time.RFC3339)
	}

	return nil
}

func (bi *BuildInfo) save(fileName string) error {
	content, err := json.MarshalIndent(bi, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal build info: %w", err)
	}

	return os.WriteFile(fileName, content, 0600)
}

func (bi *BuildInfo) loadVersion(config *Config) error { //nolint:gocyclo
	var envVersion, tagVersion, fileVersion, lastTagVersion string
	var err error

	if config.InputVersionEnvVar.EnvVar != "" {
		if envVarValue, ok := os.LookupEnv(config.InputVersionEnvVar.EnvVar); ok {
			envVersion, err = getVersionFromContent(envVarValue, config.InputVersionEnvVar.Pattern)
			if err != nil {
				return fmt.Errorf("could not get version from env var %s: %w", config.InputVersionEnvVar.EnvVar, err)
			}
		}
	}

	if config.InputVersionTag.Pattern != "" {
		if bi.GitTag != "" {
			if tagVersion, err = getVersionFromContent(bi.GitTag, config.InputVersionTag.Pattern); err != nil {
				return fmt.Errorf("failed to get version from tag: %w", err)
			}
		}

		if bi.GitLastTag != "" {
			if lastTagVersion, err = getVersionFromContent(bi.GitLastTag, config.InputVersionTag.Pattern); err != nil {
				return fmt.Errorf("failed to get version from last tag: %w", err)
			}
		}
	}

	if config.InputVersionFile.File != "" {
		if fileVersion, err = getVersionFromFile(config.InputVersionFile.File, config.InputVersionFile.Pattern); err != nil {
			return fmt.Errorf("failed to get version from file: %w", err)
		}
	}

	if fileVersion == "" {
		if err = fetchPackageManagerInfo(config.Directory, bi); err != nil {
			return fmt.Errorf("failed to fetch package manager info: %w", err)
		}
	}

	switch {
	case envVersion != "":
		bi.VersionDeclared = envVersion
		bi.Version = envVersion
	case tagVersion != "":
		bi.VersionDeclared = tagVersion
		bi.Version = tagVersion
	case fileVersion != "":
		bi.VersionDeclared = fileVersion
		bi.Version = fileVersion + "-" + bi.GitSmartRef
	case lastTagVersion != "":
		bi.Version = lastTagVersion + "-" + bi.GitSmartRef
	default:
		if bi.VersionDeclared != "" {
			bi.Version = bi.VersionDeclared + "-" + bi.GitSmartRef
		}
	}

	return nil
}
