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

	"github.com/evilsocket/islazy/tui"
	"github.com/jessevdk/go-flags"
	"github.com/lmorg/readline"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
)

// // Do - Entrypoint to all completions in Wiregost
func TabCompleter(line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {

	var suggestions []string               // Selectable suggestions
	listSuggestions := map[string]string{} // Descriptions for suggestions

	args := strings.Split(string(line), " ")         // The readline input as a []string
	last := trimSpaceLeft([]rune(args[len(args)-1])) // The last char in input

	// Detect base command automatically
	var command = detectedCommand(args, *commands.Context.Menu)

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

		// If user asks for completions with "-" or "--". (Note: This takes precedence on arguments, as it is evaluated after arguments)
		if optionsAsked(args, last, command) {
			return CompleteCommandOptions(args, last, command)
		}

		// Check environment variables
		if envVarAsked(args, last) {
			return CompleteEnvironmentVariables(line, pos)
		}

		// Propose completion for args before anything else
		if arg, yes := argumentRequired(last, args, *commands.Context.Menu, command, false); yes {
			return CompleteCommandArguments(command, arg, line, pos)
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

	switch *commands.Context.Menu {
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:
		for _, cmd := range commands.CommandsByContext() {
			// for _, cmd := range commands.MainParser.Commands() {
			if strings.HasPrefix(cmd.Name, string(last)) {
				suggestions = append(suggestions, cmd.Name[pos:]+" ")
			}
		}
	case commands.GHOST_CONTEXT:
		for _, cmd := range commands.GhostParser.Commands() {
			if strings.HasPrefix(cmd.Name, string(last)) {
				suggestions = append(suggestions, cmd.Name[pos:]+" ")
			}
		}
	}

	return string(last), suggestions, listSuggestions, readline.TabDisplayGrid
}

// CompleteCommandArguments - Completes all values for arguments to a command. Arguments here are different from command options (--option).
func CompleteCommandArguments(cmd *flags.Command, arg string, line []rune, pos int) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}
	args := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(args[len(args)-1]))

	found := commands.ArgumentByName(cmd, arg)

	// Depends first on the current menu context
	switch *commands.Context.Menu {
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:

		// This switch on the name of the argument might not be useful a lot, but we never know
		switch found.Name {
		case "Path", "OtherPath":
			// Completion might differ slightly depending on the command
			switch cmd.Name {
			case constants.Cd:
				return CompleteLocalPath(line, pos)
			case "ls", "cat":
				return completeLocalPathAndFiles(line, pos)
			case constants.ModuleUse:
				// Check for subcommand: depends on module or stack command
				if len(args) == 2 {
					return completeModulePath(line, pos)
				}
				if len(args) == 3 {
					return completeStackModulePath(line, pos)
				}
			case constants.StackPop:
				return completeStackModulePath(line, pos)
			}
		case "Option":
			switch cmd.Name {
			case constants.ModuleSetOption:
				return CompleteOptionNames(line, pos)
			}
		case "Value":
			switch cmd.Name {
			case constants.ModuleSetOption:
				// args[1] is supposed to be the option name
				return CompleteOptionValues(args[1], line, pos)
			}
		case "SessionID":
			return CompleteSessionIDs(line, pos)
		case "JobID":
			return CompleteJobIDs(line, pos)
		case "Name":
			switch cmd.Name {
			case constants.WorkspaceSwitch, constants.WorkspaceUpdate:
			case constants.ModuleParseProfile, constants.ProfilesDelete:
				return CompleteProfileNames(line, pos)
			}
		case "Server":
			return CompleteServer(line, pos)

		default: // If name is empty, return
		}

	case commands.GHOST_CONTEXT:
		switch found.Name {
		case "Path", "OtherPath", "RemotePath":
			// Completion might differ slightly depending on the command
			switch cmd.Name {
			case constants.GhostCd, constants.GhostLs, constants.GhostMkdir:
				return CompleteRemotePath(line, pos)
			case constants.GhostCat, constants.GhostDownload, constants.GhostUpload, constants.GhostRm:
				return CompleteRemotePathAndFiles(line, pos)
			}
		case "LocalPath":
			switch cmd.Name {
			case constants.GhostUpload:
				return completeLocalPathAndFiles(line, pos)
			case constants.GhostDownload:
				return CompleteLocalPath(line, pos)
			}
		case "PID":
			commands.Context.Shell.MaxTabCompleterRows = 10
			return CompleteProcesses(line, pos)
		default: // If name is empty, return
		}
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

	// If command has non-filled arguments, propose them first
	if arg, yes := argumentRequired(last, args, *commands.Context.Menu, command, true); yes {
		_, suggestions, listSuggestions, tabType = CompleteCommandArguments(command, arg, line, pos)
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
//                 case "profiles", "parse_profile":
//                         comp := &profileCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "user":
//                         comp := &userCompleter{Command: cmd}
//                         options, offset = comp.Do(ctx, line, pos)
//                 case "server":
//                         comp := &serverCompleter{Command: cmd}
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
//                 }
//
//         }
//
//         return options, offset
// }
