package core

// The core package contains basic functions needed for shell usage, such as executing commands.
// and through-client shell usage (cd, ls, etc).

import (
	"fmt"
	"os/exec"

	"github.com/evilsocket/islazy/str"
)

func ExecSilent(executable string, args []string) (string, error) {
	path, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	raw, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return "", err
	} else {
		return str.Trim(string(raw)), nil
	}
}

func Exec(executable string, args []string) (string, error) {
	out, err := ExecSilent(executable, args)
	if err != nil {
		fmt.Printf("ERROR for '%s %s': %s\n", executable, args, err)
	}
	return out, err
}

func HasBinary(executable string) bool {
	if path, err := exec.LookPath(executable); err != nil || path == "" {
		return false
	}
	return true
}
