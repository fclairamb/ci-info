package main

import (
	"regexp"
)

type BuildInfo struct {
	Version           string `json:"version"`
	VersionAuto       string
	CommitHash        string `json:"git_hash"`
	CommitHashShort   string
	CommitDate        string `json:"git_date"`
	CommitDateClean   string
	CommitBranch      string `json:"git_branch"`
	CommitBranchClean string
	CommitTag         string `json:"git_tag"`
	CommitTagOrBranch string
	BuildDate         string `json:"build_date"`
	BuildHost         string `json:"build_host"`
	BuildUser         string `json:"build_user"`
}

var reDateClean = regexp.MustCompile(`[^0-9]+`)
var reBranchClean = regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)

func (bi *BuildInfo) complete() {
	if bi.CommitDate != "" {
		date := reDateClean.ReplaceAllString(bi.CommitDate, "")
		if len(date) > 14 {
			bi.CommitDateClean = date[14:]
		} else {
			bi.CommitDateClean = date
		}
	}
	if bi.CommitBranch != "" {
		bi.CommitBranchClean = reBranchClean.ReplaceAllString(bi.CommitBranch, "-")
	}
	if bi.CommitTag != "" {
		bi.CommitTagOrBranch = bi.CommitTag
	} else if bi.CommitBranchClean != "" {
		bi.CommitTagOrBranch = bi.CommitBranchClean
	}
}
