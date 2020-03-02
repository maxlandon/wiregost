// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package completers

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/commands"
)

// AutoCompleter - The main/root autocompleter for the console.
type AutoCompleter struct {
	MenuContext *string
	Context     *commands.ShellContext
}

// Do is the completion function triggered at each line
func (ac *AutoCompleter) Do(line []rune, pos int) (options [][]rune, offset int) {

	commands := buildCommandMap(*ac.MenuContext)

	// Find commands
	var verbs []string
	for cmd := range commands {
		verbs = append(verbs, cmd)
	}

	sort.Strings(verbs)

	line = trimSpaceLeft(line[:pos])

	// Auto-complete verb
	var verbFound string
	for _, verb := range verbs {
		search := verb + " "
		if !hasPrefix(line, []rune(search)) {
			sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
			options = append(options, sLine...)
			offset = sOffset
		} else {
			verbFound = verb
			break
		}
	}
	if len(verbFound) == 0 {
		return
	}

	// Help commands need to be filtered here depending on context
	if verbFound == "help" {
		if *ac.Context.MenuContext == "agent" {
			options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
		}
	}

	// Autocomplete commands with no subcommands but variable arguments
	for _, c := range commands {
		if c.Name == "set" || c.Name == "use" || c.Name == "parse_profile" || c.Name == "help" || c.Name == "cd" {
			options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
		}
	}

	// Autocomplete subcommands
	var subFound string
	line = trimSpaceLeft(line[len(verbFound):])
	for _, comm := range commands[verbFound].SubCommands {
		search := comm + " "
		if !hasPrefix(line, []rune(search)) {
			sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
			options = append(options, sLine...)
			offset = sOffset
		} else {
			subFound = comm
			break
		}
	}

	// Option name is found, yield option completions
	if (len(subFound) == 0) && verbFound != "set" {
		return
	}
	if subFound == "interact" {
		options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
	}

	// Get command completer and yield options
	options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)

	return options, offset
}

// buildCommandMap - Creates a map of commands for completion
func buildCommandMap(ctx string) (commandMap map[string]*commands.Command) {
	commandMap = map[string]*commands.Command{}
	for _, cmd := range commands.AllContextCommands(ctx) {
		commandMap[cmd.Name] = cmd
	}
	return commandMap
}

// yieldCommandCompletions determines the type of command used and redirects to its completer
func yieldCommandCompletions(ctx *commands.ShellContext, cmd *commands.Command, line []rune, pos int) (options [][]rune, offset int) {

	switch *ctx.MenuContext {
	case "main", "module":
		switch cmd.Name {
		case "cd":
			comp := &pathCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "workspace":
			comp := &workspaceCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "hosts":
			comp := &hostCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "set":
			comp := &optionCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "use":
			comp := &moduleCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "stack":
			comp := &stackCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "profiles", "parse_profile":
			comp := &profileCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "user":
			comp := &userCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "server":
			comp := &serverCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "sessions":
			comp := &sessionCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		}

	case "agent":
		switch cmd.Name {
		case "help":
			comp := &agentHelpCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "cd":
			// Enable only if enabled in config
			if ctx.SessionPathComplete {
				comp := &implantPathCompleter{Command: cmd}
				options, offset = comp.Do(ctx, line, pos)
			}
		}

	}

	return options, offset
}

// yieldCommandCompletions determines the type of command used and redirects to its completer
func yieldOptionompletions(ctx *commands.ShellContext, cmd *commands.Command, line []rune, pos int) (options [][]rune, offset int) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	switch *ctx.MenuContext {
	case "module":
		switch cmd.Name {
		case "set":
			// If name is identified, that means option is already typed
			switch string(line) {

			case "StageImplant ", "StageConfig ":
				comp := &stagerCompleter{Command: cmd}
				options, offset = comp.Do(ctx, line, pos)

				// Default is: no options have been typed yet
			default:
				comp := &optionCompleter{Command: cmd}
				options, offset = comp.Do(ctx, line, pos)
			}
		}
	}

	return options, offset
}

// Utility functions -------------------------------------------------------------------------------------------------------------//

func trimSpaceLeft(in []rune) []rune {
	firstIndex := len(in)
	for i, r := range in {
		if unicode.IsSpace(r) == false {
			firstIndex = i
			break
		}
	}
	return in[firstIndex:]
}

func equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func hasPrefix(r, prefix []rune) bool {
	if len(r) < len(prefix) {
		return false
	}
	return equal(r[:len(prefix)], prefix)
}

func inArray(s string, array []string) bool {
	for _, item := range array {
		if s == item {
			return true
		}
	}
	return false
}

func lastString(array []string) string {
	return array[len(array)-1]
}

type argOption struct {
	Value  string
	Detail string
}

func doInternal(line []rune, pos int, lineLen int, argName []rune) (newLine [][]rune, offset int) {
	offset = lineLen
	if lineLen >= len(argName) {
		if hasPrefix(line, argName) {
			if lineLen == len(argName) {
				newLine = append(newLine, []rune{' '})
			} else {
				newLine = append(newLine, argName)
			}
			offset = offset - len(argName) - 1
		}
	} else {
		if hasPrefix(argName, line) {
			newLine = append(newLine, argName[offset:])
		}
	}
	return
}

var (
	// Info - All normal message
	Info = fmt.Sprintf("%s[-]%s ", tui.BLUE, tui.RESET)
	// Warn - Errors in parameters, notifiable events in modules/sessions
	Warn = fmt.Sprintf("%s[!]%s ", tui.YELLOW, tui.RESET)
	// Error - Error in commands, filters, modules and implants.
	Error = fmt.Sprintf("%s[!]%s ", tui.RED, tui.RESET)
	// Success - Success events
	Success = fmt.Sprintf("%s[*]%s ", tui.GREEN, tui.RESET)

	// Infof - formatted
	Infof = fmt.Sprintf("%s[-] ", tui.BLUE)
	// Warnf - formatted
	Warnf = fmt.Sprintf("%s[!] ", tui.YELLOW)
	// Errorf - formatted
	Errorf = fmt.Sprintf("%s[!] ", tui.RED)
	// Sucessf - formatted
	Sucessf = fmt.Sprintf("%s[*] ", tui.GREEN)

	//RPCError - Errors from the server
	RPCError = fmt.Sprintf("%s[RPC Error]%s ", tui.RED, tui.RESET)
	// CommandError - Command input error
	CommandError = fmt.Sprintf("%s[Command Error]%s ", tui.RED, tui.RESET)
	// DBError - Data Service error
	DBError = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)
)
