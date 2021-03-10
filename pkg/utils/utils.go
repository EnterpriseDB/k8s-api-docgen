// Package utils contains auxiliary functions for doc generation
package utils

import (
	"strings"
)

// SplitStringIntoSlice given a string and a separator, return the slice of strings after splitting with that separator
func SplitStringIntoSlice(s string, sep string) []string {
	ss := strings.Split(s, sep)
	if len(s) == 0 {
		return ss
	}
	return ss[0 : len(ss)-1]
}
