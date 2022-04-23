package main

const (
    // BUILD_VERSION is the current version of the program
    BUILD_VERSION = "{{ .Version }}"

    // BUILD_TIME is the time the program was built
    BUILD_TIME = "{{ .BuildTime }}"

    // COMMIT_SMART is the git hash of the program
    COMMIT_SMART = "{{ .CommitSmart }}"
)
