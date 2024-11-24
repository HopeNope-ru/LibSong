package utils

import "strings"

func ToVerseList(text string) []string {
	return strings.Split(text, "\\n\\n")
}
