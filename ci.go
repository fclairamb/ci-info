package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/inconshreveable/log15"
)

const true = "true"
const refBranch = "refs/heads/"
const refTags = "refs/tags/"

// CIInfoFetcher describes how we shall fetch information
type CIInfoFetcher interface {
	Detect() bool              // Detect if it's a suited fetcher
	Fetch(bi *BuildInfo) error // Fetch
	String() string            // Name of the fetcher
}

// CircleCIInfoFetcher is a fetcher for CircleCI
type CircleCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (c CircleCIInfoFetcher) Detect() bool {
	return os.Getenv("CIRCLECI") == true
}

// Fetch fetches the CI information
func (c CircleCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CIRCLE_SHA1")
	bi.CommitTag = os.Getenv("CIRCLE_TAG")
	bi.CommitBranch = os.Getenv("CIRCLE_BRANCH")

	return nil
}

func (c CircleCIInfoFetcher) String() string {
	return "circleci"
}

// GithubActionsCIInfoFetcher is a fetcher for Github Actions
// See https://docs.github.com/en/actions/learn-github-actions/environment-variables
type GithubActionsCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f GithubActionsCIInfoFetcher) Detect() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

// Fetch fetches the CI information
func (f GithubActionsCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("GITHUB_SHA")
	ref := os.Getenv("GITHUB_REF")

	if strings.HasPrefix(ref, refBranch) {
		bi.CommitBranch = ref[len(refBranch):]
	} else if strings.HasPrefix(ref, refTags) {
		bi.CommitTag = ref[len(refTags):]
	}

	return nil
}

func (f GithubActionsCIInfoFetcher) String() string {
	return "github-actions"
}

// TravisCIInfoFetcher is a fetcher for TravisCI
// See https://docs.travis-ci.com/user/environment-variables/
type TravisCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (t TravisCIInfoFetcher) Detect() bool {
	return os.Getenv("TRAVIS") == "true"
}

// Fetch fetches the CI information
func (t TravisCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("TRAVIS_COMMIT")
	bi.CommitTag = os.Getenv("TRAVIS_TAG")
	bi.CommitBranch = os.Getenv("TRAVIS_BRANCH")

	return nil
}

func (t TravisCIInfoFetcher) String() string {
	return "travis"
}

// GitLabInfoFetcher is a fetcher for GitLab CI
type GitLabInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f GitLabInfoFetcher) Detect() bool {
	return os.Getenv("GITLAB_USER_ID") != ""
}

// Fetch fetches the CI information
func (f GitLabInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CI_COMMIT_SHA")
	bi.CommitTag = os.Getenv("CI_COMMIT_TAG")
	bi.CommitBranch = os.Getenv("CI_COMMIT_REF_NAME")

	return nil
}

func (f GitLabInfoFetcher) String() string {
	return "gitlab"
}

// DroneCIInfoFetcher is a fetcher for Drone CI
// see https://docs.drone.io/pipeline/environment/reference/
type DroneCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f DroneCIInfoFetcher) Detect() bool {
	return os.Getenv("DRONE") == "true"
}

// Fetch fetches the CI information
func (f DroneCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("DRONE_COMMIT")
	bi.CommitTag = os.Getenv("DRONE_TAG")
	bi.CommitBranch = os.Getenv("DRONE_BRANCH")

	return nil
}

func (f DroneCIInfoFetcher) String() string {
	return "drone"
}

// JenkinsCIInfoFetcher is a fetcher for Jenkins CI
// see https://docs.travis-ci.com/user/environment-variables/
type JenkinsCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f JenkinsCIInfoFetcher) Detect() bool {
	return os.Getenv("JENKINS_URL") != ""
}

// Fetch fetches the CI information
func (f JenkinsCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("GIT_COMMIT")
	bi.CommitTag = os.Getenv("GIT_TAG")
	bi.CommitBranch = os.Getenv("GIT_BRANCH")

	return nil
}

func (f JenkinsCIInfoFetcher) String() string {
	return "jenkins"
}

var fetchers = []CIInfoFetcher{
	&CircleCIInfoFetcher{},
	&GithubActionsCIInfoFetcher{},
	&GitLabInfoFetcher{},
	&DroneCIInfoFetcher{},
	&TravisCIInfoFetcher{},
	&JenkinsCIInfoFetcher{},
}

func fetchCIInfo(bi *BuildInfo) error {
	for _, fetcher := range fetchers {
		if !fetcher.Detect() {
			continue
		}

		log.Info("Found CI info fetcher", "fetcher", fetcher.String())

		if err := fetcher.Fetch(bi); err != nil {
			return fmt.Errorf("failed to fetch CI info: %w", err)
		}

		break
	}

	return nil
}
