package main

import "testing"

func TestDate(t *testing.T) {
	bi := &BuildInfo{
		CommitDate: "Fri Feb  8 15:04:05 MST 2017",
	}
	bi.complete()
	if bi.CommitDateClean != "20170208" {
		t.Error("Expected 20170208, got", bi.CommitDateClean)
	}
}
