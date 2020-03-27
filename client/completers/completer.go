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
	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/commands"
)

// AutoCompleter - Handles all autocompletions and hints in Wiregost
type AutoCompleter struct {
	Context *commands.ShellContext // Passes the shell context
	Command *commands.Command      // The command currently in input
}

// Do - Entrypoint to all completions in Wiregost
func (ac *AutoCompleter) Do(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	// Get a list of words in input
	splitLine := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	// Get context commands
	cmds := buildCommandMap(*ac.Context.Menu)

	// Store Command and Subcommands
	var command *commands.Command
	for _, cmd := range cmds {
		if cmd.Name == splitLine[0] {
			command = cmd
		}
	}
	var subs []string
	if command != nil {
		for _, sub := range command.SubCommands {
			subs = append(subs, sub.Name)
		}
		sort.Strings(subs)
	}

	// Commands
	if len(splitLine) == 0 || len(splitLine) == 1 && command == nil {
		var verbs []string
		for _, cmd := range cmds {
			verbs = append(verbs, cmd.Name)
		}
		sort.Strings(verbs)

		for i := range verbs {
			if strings.HasPrefix(verbs[i], string(last)) {
				suggestions = append(suggestions, verbs[i][pos:]+" ")
			}
		}
	}

	// SubCommands
	if command != nil {
		if command.SubCommands != nil && len(command.SubCommands) != 0 {
			var verbs []string
			for _, cmd := range command.SubCommands {
				verbs = append(verbs, cmd.Name)
			}
			sort.Strings(verbs)

			for i := range verbs {
				if strings.HasPrefix(verbs[i], string(last)) {
					suggestions = append(suggestions, verbs[i][(len(last)):]+" ")
				}
			}
		}
		// If no subcommands, check Arguments
		if command.SubCommands == nil || len(command.SubCommands) == 0 {
			return ac.yieldArgsCompletion(command, line, pos)
		}
	}

	// Arguments
	if command != nil && (command.SubCommands != nil || len(command.SubCommands) != 0) && len(splitLine) > 1 && stringInSlice(splitLine[1], subs) {
		return ac.yieldArgsCompletion(command, line, pos)
	}

	return string(line[:pos]), suggestions, listSuggestions, readline.TabDisplayGrid
}

func (ac *AutoCompleter) yieldArgsCompletion(cmd *commands.Command, line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	// Completions
	var suggestions []string
	listSuggestions := map[string]string{}

	switch *ac.Context.Menu {
	case "main", "module":
		switch cmd.Name {
		case "cd":
			return completeLocalPath(cmd, line, pos)
		case "workspace":
			return completeWorkspaces(cmd, line, pos)
		}
	case "agent":
	}

	return string(line[:pos]), suggestions, listSuggestions, readline.TabDisplayGrid
}

func (ac *AutoCompleter) CommandHint(line []rune, pos int) (hint []rune) {

	splitLine := strings.Split(string(line), " ")
	line = trimSpaceLeft([]rune(splitLine[len(splitLine)-1]))

	commands := buildCommandMap(*ac.Context.Menu)

	for _, cmd := range commands {
		if cmd.Name == splitLine[0] {
			// Get main command hint
			hint = []rune(cmd.Help)

			// Commands with no subcommands but with arguments
			if len(splitLine) == 2 && splitLine[1] != "" {
				switch cmd.Name {
				case "cd":
					hint = []rune("Change directory:" + " => " + splitLine[1])
				}
			}

			// Get potential subcommands & filter hints
			var verbs []string
			for _, cmd := range cmd.SubCommands {
				verbs = append(verbs, cmd.Name)
			}
			sort.Strings(verbs)
			if len(splitLine) > 2 && stringInSlice(splitLine[1], verbs) {
				for _, sub := range cmd.SubCommands {
					if sub.Name == splitLine[1] {
						hint = []rune(sub.Help)
					}
				}

				// switch cmd.Name {
				// case "hosts":
				//         ac.Context.Shell.HintFormatting = fmt.Sprintf("%s%s", tui.YELLOW, tui.BOLD)
				//         hint = []rune("Host filters")
				// }

			}
		}
	}

	return
}

// Do is the completion function triggered at each line
// func (ac *AutoCompleter) Do(line []rune, pos int) (options [][]rune, offset int) {
//
//         commands := buildCommandMap(*ac.MenuContext)
//
//         // Find commands
//         var verbs []string
//         for cmd := range commands {
//                 verbs = append(verbs, cmd)
//         }
//
//         sort.Strings(verbs)
//
//         line = trimSpaceLeft(line[:pos])
//
//         // Auto-complete verb
//         var verbFound string
//         for _, verb := range verbs {
//                 search := verb + " "
//                 if !hasPrefix(line, []rune(search)) {
//                         sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
//                         options = append(options, sLine...)
//                         offset = sOffset
//                 } else {
//                         verbFound = verb
//                         break
//                 }
//         }
//         if len(verbFound) == 0 {
//                 return
//         }
//
//         // Help commands need to be filtered here depending on context
//         if verbFound == "help" {
//                 if *ac.Context.Menu == "agent" {
//                         options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
//                 }
//         }
//
//         // Autocomplete commands with no subcommands but variable arguments
//         for _, c := range commands {
//                 switch c.Name {
//                 case "set", "use", "parse_profile", "help", "cd", "ls", "sessions", "interact", "nmap":
//                         options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
//                 }
//         }
//
//         // Autocomplete subcommands
//         var subFound string
//         line = trimSpaceLeft(line[len(verbFound):])
//         for _, comm := range commands[verbFound].SubCommands {
//                 search := comm + " "
//                 if !hasPrefix(line, []rune(search)) {
//                         sLine, sOffset := doInternal(line, pos, len(line), []rune(search))
//                         options = append(options, sLine...)
//                         offset = sOffset
//                 } else {
//                         subFound = comm
//                         break
//                 }
//         }
//
//         // Option name is found, yield option completions
//         if (len(subFound) == 0) && verbFound != "set" {
//                 return
//         }
//         if subFound == "interact" {
//                 options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
//         }
//
//         // Get command completer and yield options
//         options, offset = yieldCommandCompletions(ac.Context, commands[verbFound], line, pos)
//
//         return options, offset
// }
//
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

	switch *ctx.Menu {
	case "main", "module":
		switch cmd.Name {
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
		case "sessions", "interact":
			comp := &sessionCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "nmap", "db_nmap":
			comp := &nmapCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		}

	case "agent":
		switch cmd.Name {
		case "help":
			comp := &agentHelpCompleter{Command: cmd}
			options, offset = comp.Do(ctx, line, pos)
		case "cd", "ls", "cat":
			// Enable only if enabled in config
			if ctx.Config.SessionPathCompletion {
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

	switch *ctx.Menu {
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
