// Wiregost - Post-Exploitation & Implant Framework
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
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/lmorg/readline"
)

// TabCompleter - Entrypoint to all tab completions in the Wiregost console.
func TabCompleter(line []rune, pos int) (lastWord string, suggestions []string, descriptions map[string]string, tabType readline.TabDisplayType) {

	// Format and sanitize input
	args, last, lastWord := FormatInput(line)

	// Detect base command automatically
	var command = detectedCommand(args, "") // add *commands.Context.Menu in the string here

	// Propose commands
	if noCommandOrEmpty(args, last, command) {

	}

	// Check environment variables
	if envVarAsked(args, last) {

	}

	// Base command is identified
	if commandFound(command) {

		// If user asks for completions with "-" / "--", show command options
		if optionsAsked(args, last, command) {

		}

		// Check environment variables again
		if envVarAsked(args, last) {

		}

		// Propose argument completion before anything, and if needed
		if _, yes := argumentRequired(last, args, "", command, false); yes { // add *commands.Context.Menu in the string here

		}

		// Then propose subcommands
		if hasSubCommands(command, args) {

		}

		// Handle subcommand if found (maybe we should rewrite this function and use it also for base command)
		if _, ok := subCommandFound(last, args, command); ok {

		}
	}

	return
}

// [ Main Completion Functions ] -----------------------------------------------------------------------------------------------------------------

// CompleteMenuCommands - Selects all commands available in a given context and returns them as suggestions
func CompleteMenuCommands(last []rune, pos int) (lastWord string, suggestions []string, descriptions map[string]string, tabType readline.TabDisplayType) {

	return string(last), suggestions, descriptions, readline.TabDisplayGrid
}

// CompleteCommandArguments - Completes all values for arguments to a command. Arguments here are different from command options (--option).
func CompleteCommandArguments(cmd *flags.Command, arg string, line []rune, pos int) (lastWord string, suggestions []string, descriptions map[string]string, tabType readline.TabDisplayType) {

	_, _, lastWord = FormatInput(line)

	return lastWord, suggestions, descriptions, readline.TabDisplayGrid
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
func HandleSubCommand(line []rune, pos int, command *flags.Command) (lastWord string, suggestions []string, descriptions map[string]string, tabType readline.TabDisplayType) {

	_, _, lastWord = FormatInput(line)

	return lastWord, suggestions, descriptions, tabType
}

// CompleteCommandOptions - Yields completion for options of a command, with various decorators
func CompleteCommandOptions(args []string, last []rune, cmd *flags.Command) (lastWord string, suggestions []string, descriptions map[string]string, tabType readline.TabDisplayType) {

	return string(last), suggestions, descriptions, readline.TabDisplayList
}

// RecursiveGroupCompletion - Handles recursive completion for nested option groups
func RecursiveGroupCompletion(args []string, last []rune, group *flags.Group) (string, []string, map[string]string, readline.TabDisplayType) {
	var suggestions []string
	listSuggestions := map[string]string{}

	return string(last), suggestions, listSuggestions, readline.TabDisplayList
}
