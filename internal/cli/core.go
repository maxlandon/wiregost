package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

// Navigation mode: Vim or Emacs.
func setModeHandler(args []string, mode bool) (new bool) {
	filter := ""
	if len(args) == 2 {
		filter = str.Trim(args[1])
	}

	switch filter {
	case "":
		if mode == true {
			println("Current mode: " + tui.Yellow("Vim"))
		} else {
			println("Current mode: " + tui.Yellow("Emacs"))
		}
		return mode
	case "vim":
		println("Switched mode: " + tui.Yellow("Vim"))
		return true
	case "emacs":
		println("Switched mode: " + tui.Yellow("Emacs"))
		return false
	}
	return true
}

// Shell Command // BROKEN !!!! NEEDS TO FIX THE args[0] (doesnt take arguments)
func shellHandler(args []string) error {
	out, err := CmdShell(args[0])
	if err == nil {
		fmt.Printf("%s\n", out)
	}
	return err
}

func CmdShell(cmd string) (string, error) {
	return Exec("/bin/sh", []string{"-c", cmd})
}

func Exec(executable string, args []string) (string, error) {
	out, err := ExecSilent(executable, args)
	if err != nil {
		fmt.Printf("ERROR for '%s %s': %s\n", executable, args, err)
	}
	return out, err
}

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

// Change directory
func changeDirHandler(args []string) string {
	raw := str.Trim(args[1])
	dir, _ := fs.Expand(raw)
	os.Chdir(dir)

	return dir
}

// Exit
func exit() {
	fmt.Println(tui.Red("[!] ") + "Quitting")
	os.Exit(0)
}
