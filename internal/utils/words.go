package utils

import (
	"regexp"
	"strings"
)

func Clear(text string) string {

	re := regexp.MustCompile(`[^\w\s]`)

	cleaned := re.ReplaceAllString(text, " ")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

	cleaned = strings.TrimSpace(cleaned)

	cleaned = strings.ToLower(cleaned)

	return cleaned
}