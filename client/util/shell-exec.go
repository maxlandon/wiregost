package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/evilsocket/islazy/str"
)

// Shell - Use the system shell transparently through the console
func Shell(args []string) error {
	out, err := Exec(args[0], args[1:])
	if err != nil {
		fmt.Printf(CommandError+"%s \n", err.Error())
		return nil
	}

	// Print output
	fmt.Println(out)

	return nil
}

// Exec - Execute a program
func Exec(executable string, args []string) (string, error) {
	path, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(path, args...)

	// Load OS environment
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}
	return str.Trim(string(out)), nil
}

// inputIsBinary - Check if first input is a system program
func inputIsBinary(args []string) bool {
	_, err := exec.LookPath(args[0])
	if err != nil {
		return false
	}
	return true
}

//
// // LoadSystemEnv - Loads all system environment variables
// func LoadSystemEnv() error {
//         env := os.Environ()
//
//         for _, kv := range env {
//                 key := strings.Split(kv, "=")[0]
//                 value := strings.Split(kv, "=")[1]
//                 SystemEnv[key] = value
//         }
//         return nil
// }
