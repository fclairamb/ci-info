package main

import "flag"

// CmdParams contains the command line parameters
type CmdParams struct {
	ConfigFile          string
	OutputBuildInfoFile string
	OutputVersionFile   string
	Version             bool
	Init                bool
}

func getParams(args []string) (*CmdParams, error) {
	fs := flag.NewFlagSet("ci-info", flag.ContinueOnError)

	params := &CmdParams{}
	fs.StringVar(&params.ConfigFile, "c", "", "config file")
	fs.BoolVar(&params.Version, "v", false, "version")
	fs.StringVar(&params.OutputBuildInfoFile, "b", "", "build info file")
	fs.StringVar(&params.OutputVersionFile, "vf", "", "version file")
	fs.BoolVar(&params.Init, "i", false, "init config file")

	return params, fs.Parse(args)
}
