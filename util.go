package main

import (
	"strings"
)

func CheckIfBot(userAgent string) bool {
	botStrings := []string{
		"bot",
		"crawl",
		"slurp",
		"spider",
		"curl",
		"wget",
	}
	userAgent = strings.ToLower(userAgent)

	for _, botStr := range botStrings {
		if strings.Contains(userAgent, botStr) {
			return true
		}
	}
	return false
}
