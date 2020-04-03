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
	"strings"
	"unicode"

	"github.com/evilsocket/islazy/tui"
	"github.com/jessevdk/go-flags"
	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/commands"
)

// // Do - Entrypoint to all completions in Wiregost
func TabCompleter(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	var suggestions []string               // Selectable suggestions
	listSuggestions := map[string]string{} // Descriptions for suggestions

	args := strings.Split(string(line), " ")         // The readline input as a []string
	last := trimSpaceLeft([]rune(args[len(args)-1])) // The last char in input

	// Detect base command automatically
	var command = detectedCommand(args)

	// Propose Commands
	if noCommandOrEmpty(args, last, command) {
		return CompleteMenuCommands(last, pos)
	}

	// Check environment variables
	if envVarAsked(args, last) {
		return CompleteEnvironmentVariables(line, pos)
	}

	// Command is identified
	if commandFound(command) {

		// Check environment variables
		if envVarAsked(args, last) {
			return CompleteEnvironmentVariables(line, pos)
		}

		// Propose completion for args before anything else
		if hasArgs(command) {
			return CompleteCommandArguments(command, line, pos)
		}

		// Then propose subcommands
		if hasSubCommands(command, args) {
			return CompleteSubCommands(args, last, command)
		}

		// Handle subcommand if found
		if sub, ok := subCommandFound(last, args, command); ok {
			return HandleSubCommand(line, pos, sub)
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

// [ Main Completion Functions ] -----------------------------------------------------------------------------------------------------------------

// CompleteMenuCommands - Selects all commands available in a given context and returns them as suggestions
func CompleteMenuCommands(last []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	cmds := commands.CommandsByContext()
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd.Name, string(last)) {
			suggestions = append(suggestions, cmd.Name[pos:]+" ")
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

// CompleteCommandArguments - Completes all values for arguments to a command. Arguments here are different from command options (--option).
func CompleteCommandArguments(cmd *flags.Command, line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	args := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(args[len(args)-1]))

	// Get commandArgs
	cmdArgs := cmd.Args()

	switch *commands.Context.Menu {
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:
		switch cmd.Name {
		case "cd":
			if len(cmdArgs) == 1 {
				arg := cmdArgs[0]
				switch arg.Name {
				case "Path":
					return completeLocalPath(line, pos)
				}
			}
			// return completeLocalPath(line, pos)

			// case "workspace":
			//         return completeWorkspaces(cmd, line, pos)
		}
	case "agent":
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

// CompleteSubCommands - Takes subcommands and gives them as suggestions
func CompleteSubCommands(args []string, last []rune, command *flags.Command) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	for _, sub := range command.Commands() {
		if strings.HasPrefix(sub.Name, string(last)) {
			suggestions = append(suggestions, sub.Name[(len(last)):]+" ")
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

// HandleSubCommand - Handles completion for subcommand options and arguments, + any option value related completion
func HandleSubCommand(line []rune, pos int, command *flags.Command) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}
	args := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(args[len(args)-1]))
	var tabType readline.TabDisplayType

	// Check environment variables
	if envVarAsked(args, last) {
		return CompleteEnvironmentVariables(line, pos)
	}

	// If command has arguments, propose them first
	if hasArgs(command) {
		_, suggestions, listSuggestions, tabType = CompleteCommandArguments(command, line, pos)
	}

	// If user asks for completions with "-" or "--". (Note: This takes precedence on arguments, as it is evaluated after arguments)
	if optionsAsked(args, last, command) {
		return CompleteCommandOptions(args, last, command)
	}

	return string(last), suggestions, listSuggestions, tabType
}

// CompleteCommandOptions - Yields completion for options of a command, with various decorators
func CompleteCommandOptions(args []string, last []rune, cmd *flags.Command) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	var group0 *flags.Group
	groups := cmd.Groups()
	if groups != nil {
		group0 = groups[0]
	}

	if group0 != nil {
		for _, opt := range group0.Options() {

			// Check if option is already set, next option if yes
			if optionNotRepeatable(opt) && optionIsAlreadySet(args, last, opt) {
				continue
			}

			if strings.HasPrefix("--"+opt.LongName, string(last)) {
				optName := "--" + opt.LongName
				suggestions = append(suggestions, optName[(len(last)):]+" ")

				var desc string
				if opt.Required {
					desc = fmt.Sprintf("%s%sR%s %s%s", tui.RED, tui.DIM, tui.RESET, tui.DIM, opt.Description)
				} else {
					desc = fmt.Sprintf("%s%sO%s %s%s", tui.GREEN, tui.DIM, tui.RESET, tui.DIM, opt.Description)
				}

				listSuggestions[optName[(len(last)):]+" "] = tui.DIM + desc + tui.RESET
			}
		}
	}
	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}

// RecursiveGroupCompletion - Handles recursive completion for nested option groups
func RecursiveGroupCompletion(args []string, last []rune, group *flags.Group) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}

// CompleteOptionValues - Yields values completion for an option requiring arguments, if they can be completed
func CompleteOptionValues(args []string, last []rune, opt *flags.Option) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}

// // Functions from the github.com/apache-monkey -------------------------------------------------------------------------------------------------------------//

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

// // yieldCommandCompletions determines the type of command used and redirects to its completer
// func yieldCommandCompletions(ctx *commands.ShellContext, cmd *commands.Command, line []rune, pos int) (options [][]rune, offset int) {
//
//         switch *ctx.Menu {
//         case "main", "module":
//                 switch cmd.Name {
//                 case "workspace":
//                         comp := &workspaceCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "hosts":
//                         comp := &hostCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "set":
//                         comp := &optionCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "use":
//                         comp := &moduleCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "stack":
//                         comp := &stackCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "profiles", "parse_profile":
//                         comp := &profileCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "user":
//                         comp := &userCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "server":
//                         comp := &serverCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "sessions", "interact":
//                         comp := &sessionCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "nmap", "db_nmap":
//                         comp := &nmapCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 }
//
//         case "agent":
//                 switch cmd.Name {
//                 case "help":
//                         comp := &agentHelpCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "cd", "ls", "cat":
//                         // Enable only if enabled in config
//                         if ctx.Config.SessionPathCompletion {
//                                 comp := &implantPathCompleter{Command: cmd}
//                                 options, offset = comp.Do(ctx, line, pos)
//                         }
//                 }
//
//         }
//
//         return options, offset
// }
//
// // yieldCommandCompletions determines the type of command used and redirects to its completer
// func yieldOptionompletions(ctx *commands.ShellContext, cmd *commands.Command, line []rune, pos int) (options [][]rune, offset int) {
//
//         args := strings.Split(string(line), " ")
//         line = trimSpaceLeft([]rune(args[len(args)-1]))
//
//         switch *ctx.Menu {
//         case "module":
//                 switch cmd.Name {
//                 case "set":
//                         // If name is identified, that means option is already typed
//                         switch string(line) {
//
//                         case "StageImplant ", "StageConfig ":
//                                 comp := &stagerCompleter{Command: cmd}
//                                 options, offset = comp.Do(ctx, line, pos)
//
//                                 // Default is: no options have been typed yet
//                         default:
//                                 comp := &optionCompleter{Command: cmd}
//                                 options, offset = comp.Do(ctx, line, pos)
//                         }
//                 }
//         }
//
//         return options, offset
// }
//

// func lastString(array []string) string {
//         return array[len(array)-1]
// }
//
// func doInternal(line []rune, pos int, lineLen int, argName []rune) (newLine [][]rune, offset int) {
//         offset = lineLen
//         if lineLen >= len(argName) {
//                 if hasPrefix(line, argName) {
//                         if lineLen == len(argName) {
//                                 newLine = append(newLine, []rune{' '})
//                         } else {
//                                 newLine = append(newLine, argName)
//                         }
//                         offset = offset - len(argName) - 1
//                 }
//         } else {
//                 if hasPrefix(argName, line) {
//                         newLine = append(newLine, argName[offset:])
//                 }
//         }
//         return
// }

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
