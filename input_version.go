package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

// ErrNoVersionInContent is returned when no version is found in the content
var ErrNoVersionInContent = errors.New("no version found in content")

func getVersionFromFile(inputFile string, pattern string) (string, error) {
	if inputFile == "" {
		return "", nil
	}

	content, err := os.ReadFile(inputFile) //nolint:gosec

	if err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	return getVersionFromContent(string(content), pattern)
}

func getVersionFromContent(content string, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("could not compile pattern: %w", err)
	}

	match := re.FindStringSubmatch(content)

	if len(match) > 1 {
		return match[1], nil
	}

	return "", ErrNoVersionInContent
}
