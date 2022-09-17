package main

import (
	log15 "github.com/inconshreveable/log15"
)

// Log is the logger
var log log15.Logger

func initLog() {
	log = log15.New()
}

// SetLoggingLevel sets the log level
func SetLoggingLevel(level string) error {
	lvl, err := log15.LvlFromString(level)
	if err != nil {
		return err
	}

	log.SetHandler(log15.LvlFilterHandler(lvl, log15.StdoutHandler))

	return nil
}
