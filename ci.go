package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/inconshreveable/log15"
)

type CIInfoFetcher interface {
	Detect() bool
	Fetch(bi *BuildInfo) error
	String() string
}

type CircleCIInfoFetcher struct{}

func (c CircleCIInfoFetcher) Detect() bool {
	return os.Getenv("CIRCLECI") == "true"
}

func (c CircleCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CIRCLE_SHA1")
	bi.CommitTag = os.Getenv("CIRCLE_TAG")
	bi.CommitBranch = os.Getenv("CIRCLE_BRANCH")

	return nil
}

func (c CircleCIInfoFetcher) String() string {
	return "circleci"
}

// https://docs.github.com/en/actions/learn-github-actions/environment-variables
type GithubActionsCIInfoFetcher struct{}

func (g GithubActionsCIInfoFetcher) Detect() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

const refBranch = "refs/heads/"
const refTags = "refs/tags/"

func (c GithubActionsCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("GITHUB_SHA")
	ref := os.Getenv("GITHUB_REF")
	if strings.HasPrefix(ref, refBranch) {
		bi.CommitBranch = ref[len(refBranch):]
	} else if strings.HasPrefix(ref, refTags) {
		bi.CommitTag = ref[len(refTags):]
	}

	return nil
}

func (c GithubActionsCIInfoFetcher) String() string {
	return "github-actions"
}

// https://docs.travis-ci.com/user/environment-variables/
type TravisCIInfoFetcher struct{}

func (t TravisCIInfoFetcher) Detect() bool {
	return os.Getenv("TRAVIS") == "true"
}

func (t TravisCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("TRAVIS_COMMIT")
	bi.CommitTag = os.Getenv("TRAVIS_TAG")
	bi.CommitBranch = os.Getenv("TRAVIS_BRANCH")

	return nil
}

func (t TravisCIInfoFetcher) String() string {
	return "travis"
}

type GitLabInfoFetcher struct{}

func (f GitLabInfoFetcher) Detect() bool {
	return os.Getenv("GITLAB_USER_ID") != ""
}

func (f GitLabInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CI_COMMIT_SHA")
	bi.CommitTag = os.Getenv("CI_COMMIT_TAG")
	bi.CommitBranch = os.Getenv("CI_COMMIT_REF_NAME")

	return nil
}

func (t GitLabInfoFetcher) String() string {
	return "gitlab"
}

// https://docs.drone.io/pipeline/environment/reference/
type DroneCIInfoFetcher struct{}

func (f DroneCIInfoFetcher) Detect() bool {
	return os.Getenv("DRONE") == "true"
}

func (f DroneCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("DRONE_COMMIT")
	bi.CommitTag = os.Getenv("DRONE_TAG")
	bi.CommitBranch = os.Getenv("DRONE_BRANCH")

	return nil
}

func (f DroneCIInfoFetcher) String() string {
	return "drone"
}

// https://docs.travis-ci.com/user/environment-variables/
type JenkinsCIInfoFetcher struct{}

func (f JenkinsCIInfoFetcher) Detect() bool {
	return os.Getenv("JENKINS_URL") != ""
}

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
			return fmt.Errorf("Failed to fetch CI info: %w", err)
		}

		break
	}

	return nil
}
