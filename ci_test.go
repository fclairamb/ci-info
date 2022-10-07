package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInfo struct {
	CISolution string
	EnvCommon  map[string]string
	EnvTag     map[string]string
	EnvBranch  map[string]string
}

func getTestsInfo() []testInfo {
	return []testInfo{
		{
			CISolution: "github-actions",
			EnvCommon: map[string]string{
				"GITHUB_ACTION": "true",
			},
			EnvTag: map[string]string{
				"GITHUB_REF": "refs/tags/v1.2.3",
			},
			EnvBranch: map[string]string{
				"GITHUB_HEAD_REF": "feature/cool-one",
				"GITHUB_REF":      "refs/heads/feature/cool-one",
			},
		},
		{
			CISolution: "circleci",
			EnvCommon: map[string]string{
				"CIRCLECI": "true",
			},
			EnvTag: map[string]string{
				"CIRCLE_TAG": "v1.2.3",
			},
			EnvBranch: map[string]string{
				"CIRCLE_BRANCH": "feature/cool-one",
			},
		},
		{
			CISolution: "travis",
			EnvCommon: map[string]string{
				"TRAVIS": "true",
			},
			EnvTag: map[string]string{
				"TRAVIS_TAG": "v1.2.3",
			},
			EnvBranch: map[string]string{
				"TRAVIS_BRANCH": "feature/cool-one",
			},
		},
		{
			CISolution: "gitlab",
			EnvCommon: map[string]string{
				"GITLAB_USER_ID": "user",
			},
			EnvTag: map[string]string{
				"CI_COMMIT_TAG": "v1.2.3",
			},
			EnvBranch: map[string]string{
				"CI_COMMIT_REF_NAME": "feature/cool-one",
			},
		},
		{
			CISolution: "jenkins",
			EnvCommon: map[string]string{
				"JENKINS_URL": "http://jenkins",
			},
			EnvTag: map[string]string{
				"GIT_TAG": "v1.2.3",
			},
			EnvBranch: map[string]string{
				"GIT_BRANCH": "feature/cool-one",
			},
		},
		{
			CISolution: "drone",
			EnvCommon: map[string]string{
				"DRONE": "true",
			},
			EnvTag: map[string]string{
				"DRONE_TAG": "v1.2.3",
			},
			EnvBranch: map[string]string{
				"DRONE_BRANCH": "feature/cool-one",
			},
		},
	}
}

func TestGithubActions(t *testing.T) {
	a := assert.New(t)

	for _, test := range getTestsInfo() {
		t.Run(test.CISolution, func(t *testing.T) {
			for k, v := range test.EnvCommon {
				t.Setenv(k, v)
			}

			t.Run("tag", func(t *testing.T) {
				for k, v := range test.EnvTag {
					t.Setenv(k, v)
				}

				bi := &BuildInfo{}

				a.NoError(fetchCISolutionInfo("", bi))

				a.Equal(test.CISolution, bi.CISolution)
				a.Equal("v1.2.3", bi.GitTag)

				a.NoError(bi.complete())
				a.Equal("v1.2.3", bi.GitRef)

				config := createDefaultConfig()
				a.NoError(bi.loadVersion(config))
				a.Equal("1.2.3", bi.Version)
			})

			t.Run("branch", func(t *testing.T) {
				for k, v := range test.EnvBranch {
					t.Setenv(k, v)
				}

				bi := &BuildInfo{}

				a.NoError(fetchCISolutionInfo("", bi))

				a.Equal(test.CISolution, bi.CISolution)
				a.Equal("feature/cool-one", bi.GitBranch)

				a.NoError(bi.complete())

				a.Equal("", bi.GitTag)
			})
		})
	}
}
