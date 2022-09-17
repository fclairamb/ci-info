package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingLevels(t *testing.T) {
	a := assert.New(t)

	t.Run("good", func(t *testing.T) {
		for _, level := range []string{"debug", "info", "warn", "error"} {
			a.NoError(SetLoggingLevel(level))
		}
	})

	t.Run("bad", func(t *testing.T) {
		a.Error(SetLoggingLevel("bad"))
	})
}
