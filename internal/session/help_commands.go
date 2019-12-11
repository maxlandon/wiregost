package session

// This file contains all help command handlers, and their registering function.

import (
	"fmt"
	"strconv"
	"strings"

	// Third Party
	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

func (s *Session) generalHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Command Categories\n")))
	fmt.Println(tui.Dim("(Type 'help CATEGORY' to show category help)"))
	fmt.Println()

	maxLen := 0
	for _, c := range commandCategories {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range commandCategories {
		fmt.Printf("  "+tui.Yellow(pad)+" : %s\n", c.Name, c.Description)
	}

	fmt.Println()
}

func (s *Session) coreHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Core Commands\n")))

	var params string
	maxLen := 0
	for _, c := range coreCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range coreCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println()
}

func (s *Session) serverHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Server Commands\n")))

	// Commands
	var params string
	maxLen := 0
	for _, c := range serverCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range serverCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))

	// Parameters
	maxLen = 0
	for _, c := range serverParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range serverParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}

	fmt.Println()
}

func (s *Session) logHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Log Commands\n")))

	// Commands
	var params string
	maxLen := 0
	for _, c := range logCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range logCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))

	// Parameters
	maxLen = 0
	for _, c := range logParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range logParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}

	fmt.Println()
}

func (s *Session) chatHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Chat Commands\n")))

	// Commands
	var params string
	maxLen := 0
	for _, c := range chatCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range chatCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println()
}

func (s *Session) workspaceHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Workspace Commands\n")))

	// Commands
	var params string
	maxLen := 0
	for _, c := range workspaceCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range workspaceCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println(tui.Bold(tui.Blue("\n  Parameters \n")))

	// Parameters
	maxLen = 0
	for _, c := range workspaceParams {
		len := len(c.Name)
		if len > maxLen {
			maxLen = len
		}
	}
	pad = "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range workspaceParams {
		dflt := tui.Dim("(default: ") + tui.Dim(c.Default) + tui.Dim(")")
		fmt.Printf("  "+tui.Yellow(pad)+" : %s %s\n", c.Name, c.Description, dflt)
	}

	fmt.Println()
}

func (s *Session) stackHelp() {
	fmt.Println(tui.Bold(tui.Blue("\n  Module Stack Commands\n")))

	// Commands
	var params string
	maxLen := 0
	for _, c := range stackCommands {
		params = strings.Join(c.Params, " ")
		len := len(c.Name + tui.Green(params))
		if len > maxLen {
			maxLen = len
		}
	}
	pad := "%" + strconv.Itoa(maxLen) + "s"

	for _, c := range stackCommands {
		params = strings.Join(c.Params, " ")
		fmt.Printf("  "+tui.Bold(pad)+" : %s\n", c.Name+" "+tui.Green(params), c.Description)
	}

	fmt.Println()
}

func (s *Session) helpHandler(args []string, sess *Session) error {
	filter := ""
	if len(args) == 2 {
		filter = str.Trim(args[1])
	}

	if filter == "" {
		s.generalHelp()
	}
	if filter == "core" {
		s.coreHelp()
	}
	if filter == "server" {
		s.serverHelp()
	}
	if filter == "log" {
		s.logHelp()
	}
	if filter == "chat" {
		s.chatHelp()
	}
	if filter == "workspace" {
		s.workspaceHelp()
	}
	if filter == "stack" {
		s.stackHelp()
	}
	return nil
}

func (s *Session) registerHelpHandlers() {

	s.addHandler(NewCommandHandler("help",
		"^(help|\\?)(.*)$",
		"List available commands or show module specific help if no module name is provided.",
		s.helpHandler),
		readline.PcItem("help", readline.PcItemDynamic(func(prefix string) []string {
			prefix = str.Trim(prefix[4:])
			sets := []string{""}
			for _, m := range commandCategories {
				if prefix == "" || strings.HasPrefix(m.Name, prefix) {
					sets = append(sets, m.Name)
				}
			}
			return sets
		})))
}
