package session

import (
	// Standard
	"fmt"
	"os"
	"os/exec"
	"strings"

	// 3rd party
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

// Function used for description paragraphs
func wrap(text string, lineWidth int) (wrapped string) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}

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
	out, err := cmdShell(args[0])
	if err == nil {
		fmt.Printf("%s\n", out)
	}
	return err
}

func cmdShell(cmd string) (string, error) {
	return Exec("/bin/sh", []string{"-c", cmd})
}

// Exec needs to be exported because of os/exec package conflict.
func Exec(executable string, args []string) (string, error) {
	out, err := execSilent(executable, args)
	if err != nil {
		fmt.Printf("ERROR for '%s %s': %s\n", executable, args, err)
	}
	return out, err
}

func execSilent(executable string, args []string) (string, error) {
	path, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	raw, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return str.Trim(string(raw)), nil
}

// Change directory
func changeDirHandler(args []string) string {
	raw := str.Trim(args[1])
	dir, _ := fs.Expand(raw)
	os.Chdir(dir)

	return dir
}

// Set parameter option handler.
func (s *Session) setOption(cmd []string) {
	// Keep values locally for workspace
	if strings.HasPrefix(cmd[1], "workspace") {
		s.Env[cmd[1]] = strings.Join(cmd[2:], " ")
		fmt.Println()
		fmt.Printf("[-] %s%s%s set to %s%s%s", tui.YELLOW, cmd[1], tui.RESET, tui.YELLOW, s.Env[cmd[1]], tui.RESET)
		fmt.Println()
	}
	// Keep values locally for endpoint
	if strings.HasPrefix(cmd[1], "endpoint") {
		s.Env[cmd[1]] = strings.Join(cmd[2:], " ")
		fmt.Println()
		fmt.Printf("[-] %s%s%s set to %s%s%s", tui.YELLOW, cmd[1], tui.RESET, tui.YELLOW, s.Env[cmd[1]], tui.RESET)
		fmt.Println()
	}
	// Keep values locally for server. BUT
	// Loaded automatically from server when switching workspace.
	// Can be modified here and sent back when respawning a server.
	// ULTIMATELY WE SHOULD ADD CONTROLS AND/OR WARNINGS SO THAT SERVER IS NOT RELOADED WITH DIFFERENT PARAMETERS
	// THAT WOULD MAKE IT UNUSABLE BY ALREADY REGISTERED/TO-BE-REGISTERED AGENTS.
	if strings.HasPrefix(cmd[1], "server") {
		s.Env[cmd[1]] = strings.Join(cmd[2:], " ")
		fmt.Println()
		fmt.Printf("[-] %s%s%s set to %s%s%s", tui.YELLOW, cmd[1], tui.RESET, tui.YELLOW, s.Env[cmd[1]], tui.RESET)
		fmt.Println()
	}
}

func (s *Session) getOption(cmd []string) {
	var value string
	if v, ok := s.Env[cmd[1]]; ok {
		value = v
	}
	fmt.Println()
	fmt.Printf(" %s%s%s => %s%s%s \n", tui.YELLOW, cmd[1], tui.RESET, tui.YELLOW, value, tui.RESET)
}

// Exit
func exit() {
	fmt.Println(tui.Red("[!] ") + "Quitting")
	os.Exit(0)
}
