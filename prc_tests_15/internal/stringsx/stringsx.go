package stringsx

import (
	"strings"
	"unicode/utf8"
)

func Clip(s string, max int) string {
	if max < 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= max {
		return s
	}

	return string(runes[:max])
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	firstRun, size := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(firstRun)) + s[size:]
}
