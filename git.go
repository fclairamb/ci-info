package main

import (
	"os/exec"
	"strings"

	log "github.com/inconshreveable/log15"
)

func fetchGitInfoWithCmd(info *BuildInfo) error {
	var gitCommands = []gitInfoFetch{
		{&info.CommitTag, []string{"git", "tag", "--points-at", "HEAD"}},
		{&info.CommitHash, []string{"git", "rev-parse", "HEAD"}},
		{&info.CommitBranch, []string{"git", "branch", "--show-current"}},
		{&info.CommitDate, []string{"git", "show", "--quiet", "--format=%ci"}},
	}

	for _, gifc := range gitCommands {
		if *gifc.info != "" {
			continue
		}

		log.Debug("Fetching git info", "command", gifc.command)
		cmd := exec.Command(gifc.command[0], gifc.command[1:]...)
		if outBin, err := cmd.CombinedOutput(); err != nil {
			return err
		} else {
			out := strings.TrimRight(string(outBin), "\n")
			log.Debug("Fetched git info", "command", gifc.command, "output", out)
			*gifc.info = out
		}
	}

	return nil
}
