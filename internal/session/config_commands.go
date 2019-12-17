package session

// This file contains all configuration command handlers, and their
// regitering function.

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

func (s *Session) configHandler(args []string, sess *Session) error {
	var err error
	config := s.Config.ExportConfig()

	fmt.Println(tui.Bold(tui.Blue("\n  Configuration\n")))
	// Directories
	fmt.Println(tui.Blue("Directories"))
	directories := []string{
		"User directory",
		"Logs directory",
		"Personal modules",
		"Generated payloads",
		"Console resources",
		"Workspace data",
		"Export files",
	}

	values := reflect.ValueOf(config)
	num := len(directories)

	maxLen := 0
	for i := 0; i < num; i++ {
		field := directories[i]
		len := len(field)
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"
	for i := 0; i < num; i++ {
		value := values.Field(i)
		fmt.Printf("  "+tui.Yellow(pad)+" : %s \n", directories[i], value)
	}
	fmt.Println()

	// Files
	fmt.Println(tui.Blue("Files"))
	files := []string{
		"User Configuration",
		"History",
		"Global Variables",
	}

	values = reflect.ValueOf(config)
	num = len(files) + 7

	maxLen = 0
	for i := 0; i < 3; i++ {
		field := files[i]
		len := len(field)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for i := 7; i < num; i++ {
		value := values.Field(i)
		fmt.Printf("  "+tui.Yellow(pad)+" : %s \n", files[i-7], value)
	}
	fmt.Println()
	return err
}

func (s *Session) registerConfigHandlers() {
	// Core Configuration
	s.addHandler(NewCommandHandler("config",
		"^config$",
		"Show configuration variables.",
		s.configHandler),
		readline.PcItem("config"))
}
