# CI Info

This tool helps extract CI info and embed them in a resulting code (whether it's compiled or not).

## Why
Adding the CI info useful to identify what was used to build any app. Yet, doing it properly is boring.

## How it works
You provide a template file and it will take care of writing the final file with the build information.

This makes it completely language-agnostic.

## Very basic usage
You can also completely skip the templating system:
```zsh
% ci-info -b build.json -vf version.txt

% cat build.json  
{
  "ci_info_version": "0.1.0-feature-config-change-6ea1772",
  "version": "0.1.0-feature-config-change-6ea1772",
  "git_hash": "6ea17722aa995a6c69c67e833d3c5abee463f7da",
  "git_date": "2022-09-10 22:51:17 +0200",
  "git_branch": "feature/config-change",
  "build_date": "2022-09-10-2126",
  "build_host": "MBPdeFlorent2",
  "build_user": "florent"
}

% cat version.txt
0.1.0-feature-config-change-6ea1772
```

## Supported CI

The most popular continuous integration services are supported.

- [CircleCI](https://circleci.com/)
- [GitHub Actions](https://github.com/features/actions)
- [Gitlab CI](https://docs.gitlab.com/ee/ci/)
- [Drone CI](https://drone.io/)
- [Travis CI](https://travis-ci.org/)
- [Jenkins](https://jenkins.io/)

## Sample config file
The `.ci-info.json` looks like this:
```json
{
  "$schema": "https://raw.githubusercontent.com/fclairamb/ci-info/main/config-schema.json",
  "version_input_file": {
    "file": "README.md",
    "pattern": "Version: ([0-9.]+)\n"
  },
  "version_input_tag": {
    "pattern": "^v?([0-9.]+)$"
  },
  "version_input_env_var": {
    "env_var": "VERSION",
    "pattern": "^([0-9.]+)$"
  },
  "templates": [{
    "input_file": "build.go.tpl",
    "output_file": "build.go"
  }, {
    "input_content": "{{ .Version }}",
    "output_file": "version.txt"
  }],
  "build_info_file": "build.json"
}
```
## Possible template arguments
| Argument | Sample value | Description |
| -------- | ------------ | ----------- |
| `{{ .Version }}` | `0.1.0-fix-pr-check-f96a756` | The automatically generated version. This is mix of the declared one and the current GIT info. |
| `{{ .CommitHash }}` | `f96a75638b0e1767f969e23f383f4bc75c0e6ba0` | The current GIT commit |
| `{{ .CommitHashShort }}` | `f96a756` | Short version of a hash |
| `{{ .CommitDate }}` | `2022-04-23 23:52:13 +0200` | The commit's date |
| `{{ .CommitDateClean }}` | `2022-04-23-2157` | The commit's date in a clean format |
| `{{ .CommitBranch }}` | `fix/pr-check` | The current branch |
| `{{ .CommitBranchClean }}` | `fix-pr-check` | The commit branch without special chars |
| `{{ .CommitTag }}` | `v0.1.0` | The current GIT tag |
| `{{ .CommitRef }}` | `v0.1.0` | The current GIT tag or branch |
| `{{ .CommitSmart }}` | `fix-pr-check-f96a756` | The current GIT commit described by tag, otherwise branch + hash, otherwise hash |
| `{{ .BuildDate }}` | `2022-04-23-2210` | The build time |
| `{{ .BuildHost }}` | `build-server` | The build host |
| `{{ .BuildUser }}` | `runner` | The build user |
| `{{ .CISolution }}` | `circleci` | The CI solution |
| `{{ .CIBuildNumber }}` | `123` | The CI build number |
| `{{ .PackageManager }}` | `npm` | The package manager |

# Run it
## With a local binary
```sh
curl https://github.com/fclairamb/ci-info/releases/download/v0.1.10/ci-info_0.1.10_darwin_arm64.tar.gz |tar -x
./ci-info
```

## with a container
```sh
docker run --rm -ti -v $(pwd):/work fclairamb/ci-info
```