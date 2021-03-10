package utils

import (
	"os/exec"
	"strings"
)

// Given a string and a separator, return the slice of strings after splitting with that separator
func SplitStringIntoSlice(s string, sep string) []string {
	ss := strings.Split(s, sep)
	if len(s) == 0 {
		return ss
	}
	return ss[0 : len(ss)-1]
}

// Perform an ls bash command and return a slice of strings containing file names.
// e.g. if path is "." and pattern is "*types.go", then compose "ls ./*types.go" and the result would be []{"utils.go"}
func GetFilenames(path string, pattern string) ([]string, error) {
	lsCmd := exec.Command("bash", "-c", "ls "+path+`/*`+pattern+"*")
	lsOut, err := lsCmd.Output()
	if err != nil {
		return nil, err
	}
	return SplitStringIntoSlice(string(lsOut), "\n"), nil
}
