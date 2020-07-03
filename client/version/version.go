package version

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// Version - The semantic version in string form
	Version string
	// GitCommit - The commit id at compile time
	GitCommit string
	// GitDirty - Was the commit dirty at compile time
	GitDirty string
	// CompiledAt - When was this binary compiled
	CompiledAt string
)

// SemanticVersion - Get the structured sematic version
func SemanticVersion() []int {
	semVer := []int{}
	for _, part := range strings.Split(Version, ".") {
		number, _ := strconv.Atoi(part)
		semVer = append(semVer, number)
	}
	return semVer
}

// Compiled - Get time this binary was compiled
func Compiled() (time.Time, error) {
	compiled, err := strconv.ParseInt(CompiledAt, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(compiled, 0), nil
}

// FullVersion - Full version string
func FullVersion() string {
	ver := fmt.Sprintf("%s", Version)
	compiled, err := Compiled()
	if err != nil {
		ver += fmt.Sprintf(" - Compiled %s", compiled.String())
	}
	if GitCommit != "" {
		ver += fmt.Sprintf(" - %s", GitCommit)
	}
	return ver
}
