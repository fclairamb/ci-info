At [habx](https://github.com/habx/):
- The current version is always specified by the `version` property of a `package.json` file, even if the project has nothing to do with node.js.
- We produce two files:
  - A `build.json` file similar to the one produced by ci-info
  - A `version.txt` file containing a version in the form of `0.1.0` / `0.1.0-feature-config-change-6ea1772`
