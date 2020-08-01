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
	"github.com/maxlandon/readline"
)

// TabCompleter - Entrypoint to all tab completions in the Wiregost console.
func TabCompleter(line []rune, pos int) (lastWord string, completions []*readline.CompletionGroup) {

	// Format and sanitize input
	args, last, lastWord := FormatInput(line)

	// Detect base command automatically
	var command = detectedCommand(args, "") // add *commands.Context.Menu in the string here

	// Propose commands
	if noCommandOrEmpty(args, last, command) {
		return CompleteMenuCommands(last, pos)
	}

	// Check environment variables
	if envVarAsked(args, last) {

	}

	// Base command has been identified
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

	// -------------------- IMPORTANT ------------------------
	// WE NEED TO PASS A DEEP COPY OF THE OBJECTS: OTHERWISE THE COMPLETION SEARCH FUNCTION WILL MESS UP WITH THEM.

	return
}

// [ Main Completion Functions ] -----------------------------------------------------------------------------------------------------------------

// CompleteMenuCommands - Selects all commands available in a given context and returns them as suggestions
// Many categories, all from command parsers.
func CompleteMenuCommands(last []rune, pos int) (lastWord string, completions []*readline.CompletionGroup) {

	return
}

// CompleteCommandArguments - Completes all values for arguments to a command. Arguments here are different from command options (--option).
// Many categories, from multiple sources in multiple contexts
func CompleteCommandArguments(cmd *flags.Command, arg string, line []rune, pos int) (lastWord string, completions []*readline.CompletionGroup) {

	_, _, lastWord = FormatInput(line)

	return
}

// CompleteSubCommands - Takes subcommands and gives them as suggestions
// One category, from one source (a parent command)
func CompleteSubCommands(args []string, last []rune, command *flags.Command) (lastWord string, completions []*readline.CompletionGroup) {

	for _, sub := range command.Commands() {
		if strings.HasPrefix(sub.Name, string(last)) {
			// suggestions = append(suggestions, sub.Name[(len(last)):]+" ")
		}
	}

	return
}

// HandleSubCommand - Handles completion for subcommand options and arguments, + any option value related completion
// Many categories, from many sources: this function calls the same functions as the ones previously called for completing its parent command.
func HandleSubCommand(line []rune, pos int, command *flags.Command) (lastWord string, completions []*readline.CompletionGroup) {

	_, _, lastWord = FormatInput(line)

	return
}

// CompleteCommandOptions - Yields completion for options of a command, with various decorators
// Many categories, from one source (a command)
func CompleteCommandOptions(args []string, last []rune, cmd *flags.Command) (lastWord string, completions []*readline.CompletionGroup) {

	return
}

// RecursiveGroupCompletion - Handles recursive completion for nested option groups
// Many categories, one source (a command's root option group). Called by the function just above.
func RecursiveGroupCompletion(args []string, last []rune, group *flags.Group) (lastWord string, completions []*readline.CompletionGroup) {

	return
}
