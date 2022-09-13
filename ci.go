package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/inconshreveable/log15"
)

const sTrue = "true"
const refBranch = "refs/heads/"
const refTags = "refs/tags/"

var errCouldNotFindVersion = errors.New("could not find version")

var ciSolutionsFetchers = []CIInfoFetcher{
	&circleCIInfoFetcher{},
	&githubActionsCIInfoFetcher{},
	&gitLabInfoFetcher{},
	&droneCIInfoFetcher{},
	&travisCIInfoFetcher{},
	&jenkinsCIInfoFetcher{},
}

var packageManagerFetchers = []CIInfoFetcher{
	&npmInfoFetcher{},
	&gradleInfoFetcher{},
	&mavenInfoFetcher{},
	&nugetInfoFetcher{},
}

func fetchCISolutionInfo(bi *BuildInfo) error {
	return fetchCIInfo(bi, ciSolutionsFetchers, &bi.CISolution)
}

func fetchPackageManagerInfo(bi *BuildInfo) error {
	return fetchCIInfo(bi, packageManagerFetchers, &bi.PackageManager)
}

func fetchCIInfo(bi *BuildInfo, fetchers []CIInfoFetcher, target *string) error {
	for _, fetcher := range fetchers {
		if !fetcher.Detect() {
			continue
		}

		log.Info("Found CI info fetcher", "fetcher", fetcher.String())

		if err := fetcher.Fetch(bi); err != nil {
			return fmt.Errorf("failed to fetch CI info: %w", err)
		}

		*target = fetcher.String()

		break
	}

	return nil
}

// CIInfoFetcher describes how we shall fetch information
type CIInfoFetcher interface {
	Detect() bool              // Detect if it's a suited fetcher
	Fetch(bi *BuildInfo) error // Fetch
	String() string            // Name of the fetcher
}

// circleCIInfoFetcher is a fetcher for CircleCI
// see https://circleci.com/docs/variables
type circleCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (c circleCIInfoFetcher) Detect() bool {
	return os.Getenv("CIRCLECI") == sTrue
}

// Fetch fetches the CI information
func (c circleCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CIRCLE_SHA1")
	bi.CommitTag = os.Getenv("CIRCLE_TAG")
	bi.CommitBranch = os.Getenv("CIRCLE_BRANCH")
	bi.CIBuildNumber = os.Getenv("CIRCLE_BUILD_NUM")

	return nil
}

func (c circleCIInfoFetcher) String() string {
	return "circleci"
}

// githubActionsCIInfoFetcher is a fetcher for Github Actions
// See https://docs.github.com/en/actions/learn-github-actions/environment-variables
type githubActionsCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f githubActionsCIInfoFetcher) Detect() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

// Fetch fetches the CI information
func (f githubActionsCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("GITHUB_SHA")
	bi.CIBuildNumber = os.Getenv("GITHUB_RUN_ID")
	ref := os.Getenv("GITHUB_REF")

	if strings.HasPrefix(ref, refBranch) {
		bi.CommitBranch = ref[len(refBranch):]
	} else if strings.HasPrefix(ref, refTags) {
		bi.CommitTag = ref[len(refTags):]
	}

	return nil
}

func (f githubActionsCIInfoFetcher) String() string {
	return "github-actions"
}

// travisCIInfoFetcher is a fetcher for TravisCI
// See https://docs.travis-ci.com/user/environment-variables/
type travisCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (t travisCIInfoFetcher) Detect() bool {
	return os.Getenv("TRAVIS") == "true"
}

// Fetch fetches the CI information
func (t travisCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("TRAVIS_COMMIT")
	bi.CommitTag = os.Getenv("TRAVIS_TAG")
	bi.CommitBranch = os.Getenv("TRAVIS_BRANCH")
	bi.CIBuildNumber = os.Getenv("TRAVIS_BUILD_NUMBER")

	return nil
}

func (t travisCIInfoFetcher) String() string {
	return "travis"
}

// gitLabInfoFetcher is a fetcher for GitLab CI
type gitLabInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f gitLabInfoFetcher) Detect() bool {
	return os.Getenv("GITLAB_USER_ID") != ""
}

// Fetch fetches the CI information
func (f gitLabInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("CI_COMMIT_SHA")
	bi.CommitTag = os.Getenv("CI_COMMIT_TAG")
	bi.CommitBranch = os.Getenv("CI_COMMIT_REF_NAME")
	bi.CIBuildNumber = os.Getenv("CI_PIPELINE_ID")

	return nil
}

func (f gitLabInfoFetcher) String() string {
	return "gitlab"
}

// droneCIInfoFetcher is a fetcher for Drone CI
// see https://docs.drone.io/pipeline/environment/reference/
type droneCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f droneCIInfoFetcher) Detect() bool {
	return os.Getenv("DRONE") == "true"
}

// Fetch fetches the CI information
func (f droneCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("DRONE_COMMIT")
	bi.CommitTag = os.Getenv("DRONE_TAG")
	bi.CommitBranch = os.Getenv("DRONE_BRANCH")
	bi.CIBuildNumber = os.Getenv("DRONE_BUILD_NUMBER")

	return nil
}

func (f droneCIInfoFetcher) String() string {
	return "drone"
}

// jenkinsCIInfoFetcher is a fetcher for Jenkins CI
// see https://docs.travis-ci.com/user/environment-variables/
type jenkinsCIInfoFetcher struct{}

// Detect if it's a suited fetcher
func (f jenkinsCIInfoFetcher) Detect() bool {
	return os.Getenv("JENKINS_URL") != ""
}

// Fetch fetches the CI information
func (f jenkinsCIInfoFetcher) Fetch(bi *BuildInfo) error {
	bi.CommitHash = os.Getenv("GIT_COMMIT")
	bi.CommitTag = os.Getenv("GIT_TAG")
	bi.CommitBranch = os.Getenv("GIT_BRANCH")
	bi.CIBuildNumber = os.Getenv("BUILD_NUMBER")

	return nil
}

func (f jenkinsCIInfoFetcher) String() string {
	return "jenkins"
}

type npmInfoFetcher struct{}

func (f npmInfoFetcher) Detect() bool {
	st, err := os.Stat("package.json")
	if err != nil {
		return false
	}

	return !st.IsDir()
}

// Fetch parses a package.json file and retrieve the version property
func (f npmInfoFetcher) Fetch(bi *BuildInfo) error {
	b, err := os.ReadFile("package.json")
	if err != nil {
		return err
	}

	var pkg struct {
		Version string `json:"version"`
	}

	if err := json.Unmarshal(b, &pkg); err != nil {
		return err
	}

	bi.VersionDeclared = pkg.Version

	return nil
}

func (f npmInfoFetcher) String() string {
	return "npm"
}

type gradleInfoFetcher struct{}

func (f gradleInfoFetcher) Detect() bool {
	st, err := os.Stat("build.gradle")
	if err != nil {
		return false
	}

	return !st.IsDir()
}

// Fetch parses a build.gradle file and retrieve the version property
// Status: Completely broken logic
func (f gradleInfoFetcher) Fetch(bi *BuildInfo) error {
	b, err := os.ReadFile("build.gradle")
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`version\s*=\s*['"](.*)['"]`)
	matches := re.FindStringSubmatch(string(b))

	if len(matches) != 2 {
		return errCouldNotFindVersion
	}

	bi.VersionDeclared = matches[1]

	return nil
}

func (f gradleInfoFetcher) String() string {
	return "gradle"
}

type mavenInfoFetcher struct{}

func (f mavenInfoFetcher) Detect() bool {
	st, err := os.Stat("pom.xml")
	if err != nil {
		return false
	}

	return !st.IsDir()
}

// Fetch parses a pom.xml file and retrieve the version property
// Status: Should work
func (f mavenInfoFetcher) Fetch(bi *BuildInfo) error {
	b, err := os.ReadFile("pom.xml")
	if err != nil {
		return err
	}

	var pom struct {
		Version string `xml:"version"`
	}

	if err := xml.Unmarshal(b, &pom); err != nil {
		return err
	}

	bi.VersionDeclared = pom.Version

	return nil
}

func (f mavenInfoFetcher) String() string {
	return "maven"
}

type nugetInfoFetcher struct{}

func (f nugetInfoFetcher) Detect() bool {
	st, err := os.Stat("*.csproj")
	if err != nil {
		return false
	}

	return !st.IsDir()
}

// Fetch parses a *.csproj file and retrieve the version property
// Status: Should work
func (f nugetInfoFetcher) Fetch(bi *BuildInfo) error {
	files, err := filepath.Glob("*.csproj")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("unable to find any csproj file: %w", errCouldNotFindVersion)
	}

	b, err := os.ReadFile(files[0])
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`<Version>(.*)</Version>`)
	matches := re.FindStringSubmatch(string(b))

	if len(matches) != 2 {
		return fmt.Errorf("unable to find version in %s: %w", files[0], errCouldNotFindVersion)
	}

	bi.Version = matches[1]

	return nil
}

func (f nugetInfoFetcher) String() string {
	return "nuget"
}
