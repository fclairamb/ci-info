package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func getVersionFromFile(inputFile string, pattern string) (string, error) {
	if inputFile == "" {
		return "", nil
	}
	content, err := ioutil.ReadFile(inputFile)
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

	match := re.FindStringSubmatch(string(content))

	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", fmt.Errorf("could not find version in content")
	}
}
