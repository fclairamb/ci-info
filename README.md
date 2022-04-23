# CI Info

This tool helps extract CI info.

Version: 0.1.0
## Why
Adding the CI info useful to identify what was used to build any app. Yet, doing it properly is boring.

## How
Your provide a template file and it will take care of writing the final file with the build information.

This makes it completely language agnostic.

## Supported CI

The most popular continuous integration services are suported.

- [CircleCI](https://circleci.com/)
- [GitHub Actions](https://github.com/features/actions)
- [Gitlab CI](https://docs.gitlab.com/ee/ci/)
- [Drone CI](https://drone.io/)
- [Travis CI](https://travis-ci.org/)
- [Jenkins](https://jenkins.io/)

## Possible template arguments
| Argument | Description |
| -------- | ----------- |
| {{ .Version }} | The automatically generated version. This is mix of the declared one and the current GIT info. |
| {{ .CommitHash }} | The current GIT commit |
| {{ .CommitHashShort }} | The current GIT branch |
| {{ .CommitDate }} | The current GIT tag |
| {{ .CommitDateClean }} | The current GIT status |
| {{ .CommitBranchClean }} | The commit branch without special chars |
| {{ .CommitTag }} | The current GIT tag |
| {{ .CommitRef }} | The current GIT tag or branch |
| {{ .CommitSmart }} | The current GIT commit described by tag, otherwise branch + hash, otherwise hash |
| {{ .BuildTime }} | The build time |
| {{ .BuildHost }} | The build host |
| {{ .BuildUser }} | The build user |
