package strn

import (
	"unicode/utf8"
)

//goland:noinspection GoUnusedConst
const (
	EmptyStr = ""
)

func Count(s string) int {
	return utf8.RuneCountInString(s)
}

func Sub(s string, start, end int) string {
	runes := []rune(s)
	return string(runes[start:end])
}

//goland:noinspection GoUnusedExportedFunction
func Truncate(s string, maxCount int) string {
	if l := Count(s); l > maxCount {
		return Sub(s, 0, maxCount)
	}
	return s
}
