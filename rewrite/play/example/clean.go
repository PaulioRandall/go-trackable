package example

import (
//"github.com/PaulioRandall/trackable/rewrite"
)

const (
	parseDayFmt = "13 Jan"
)

func cleanData(data [][]string) ([][]string, error) {
	panic("clean.cleanData")
}

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}
