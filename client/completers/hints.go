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
	"reflect"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/jessevdk/go-flags"
	"github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
)

var (
	// Hint signs
	menuCatHint    = tui.RESET + tui.DIM + tui.BOLD + " Menu  " + tui.RESET                              // Dim
	envCatHint     = tui.RESET + tui.GREEN + tui.BOLD + " Env  " + tui.RESET + tui.DIM + tui.GREEN       // Green
	commandCatHint = "\033[38;5;223m" + tui.BOLD + " Command  " + tui.RESET + tui.DIM + "\033[38;5;223m" // Cream
	ExeCatHint     = tui.RESET + tui.DIM + tui.BOLD + " Shell " + tui.RESET + tui.DIM                    // Dim
	optionCatHint  = "\033[38;5;222m" + tui.BOLD + " Options  " + tui.RESET + tui.DIM + "\033[38;5;222m" // Cream-Yellow
	valueCatHint   = "\033[38;5;217m" + tui.BOLD + " Value  " + tui.RESET + tui.DIM + "\033[38;5;217m"   // Pink-Cream
	argCatHint     = "\033[38;5;217m" + tui.BOLD + " Arg  " + tui.RESET + tui.DIM + "\033[38;5;217m"     // Pink-Cream
)

// HintText - Main function for displaying command/argument/option hints
func HintText(line []rune, pos int) (hint []rune) {

	// Curated input
	args := strings.Split(string(line), " ")
	last := trimSpaceLeft([]rune(args[len(args)-1]))

	// Detect base command automatically
	var command = detectedCommand(args)

	// Menu hints
	if noCommandOrEmpty(args, last, command) {
		hint = MenuHint(args, last)
	}

	// Environment variables
	if envVarAsked(args, last) {
		return envVarHint(args, last)
	}

	// Command Hint
	if commandFound(command) {
		// Command hint by default
		hint = []rune(commandCatHint + command.ShortDescription)

		// Environment variables
		if envVarAsked(args, last) {
			return envVarHint(args, last)
		}

		// If command has args, hint for args
		if arg, yes := argumentRequired(last, args, command, false); yes {
			hint = []rune(CommandArgumentHints(args, last, command, arg))
			return
		}

		// Brief subcommand hint
		if lastIsSubCommand(last, command) {
			hint = []rune(commandCatHint + command.Find(string(last)).ShortDescription)
			return
		}

		// Handle subcommand if found
		if sub, ok := subCommandFound(last, args, command); ok {
			return HandleSubcommandHints(args, last, sub)
		}
	}

	// Handle special commands
	if isSpecialCommand(args, command) {
		return handleSpecialCommands(args, line)
	}

	if commandFoundInPath(args[0]) {
		hint = []rune(ExeCatHint + util.ParseSummary(util.GetManPages(args[0])))
	}

	return
}

// HandleSubcommandHints - Handles hints for a subcommand and its arguments, options, etc.
func HandleSubcommandHints(args []string, last []rune, command *flags.Command) (hint []rune) {

	// If command has args, hint for args
	if arg, yes := argumentRequired(last, args, command, true); yes {
		hint = []rune(CommandArgumentHints(args, last, command, arg))
		return
	}

	// Environment variables
	if envVarAsked(args, last) {
		hint = envVarHint(args, last)
	}

	// If the last word in input is an option --name, yield argument hint if needed
	if len(command.Groups()) != 0 {
		if opt, yes := optionArgRequired(args, last, command.Groups()[0]); yes {
			hint = OptionArgumentHint(args, last, opt)
		}
	}

	// If user asks for completions with "-" or "--". (Note: This takes precedence on any argument hints, as it is evaluated after them)
	if optionsAsked(args, last, command) {
		return OptionHints(args, last, command)
	}

	return
}

// ArgumentHints - Yields hints for arguments to commands if they have some
func CommandArgumentHints(args []string, last []rune, command *flags.Command, arg string) (hint []rune) {

	found := commands.ArgumentByName(command, arg)

	// Base Hint is just a description of the command argument
	hint = []rune(argCatHint + found.Description)

	switch *commands.Context.Menu {
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:

		switch found.Name {
		case "Value":
			switch command.Name {
			case constants.ModuleSetOption:
				// args[1] is supposed to be the option name
				hint = ModuleOptionHints(args[1])
			}

		default: // If name is empty, return
		}

	case commands.GHOST_CONTEXT:
	}

	return
}

func ModuleOptionHints(opt string) (hint []rune) {
	hint = []rune(valueCatHint + commands.Context.Module.Options[opt].Description)
	return
}

// OptionHints - Yields hints for proposed options lists/groups
func OptionHints(args []string, last []rune, command *flags.Command) (hint []rune) {

	groups := command.Groups()
	if groups != nil {
		group0 := groups[0]
		hint = []rune(optionCatHint + group0.ShortDescription)
	}

	return
}

// OptionArgumentHint - Yields hints for arguments to an option (generally the last word in input)
func OptionArgumentHint(args []string, last []rune, opt *flags.Option) (hint []rune) {
	return []rune(valueCatHint + opt.Description)
}

// getRemainingArgs - Filters the input slice from commands and detected option:value pairs, and returns args
func getRemainingArgs(args []string, last []rune, command *flags.Command) (remain []string) {

	var input []string
	// Clean subcommand name
	if args[0] == command.Name && len(args) >= 2 {
		input = args[1:]
	} else if len(args) == 1 {
		input = args
	}

	// For each each argument
	for i := 0; i < len(input); i++ {
		// Check option prefix
		if strings.HasPrefix(input[i], "-") || strings.HasPrefix(input[i], "--") {
			// Clean it
			cur := strings.TrimPrefix(input[i], "--")
			cur = strings.TrimPrefix(cur, "-")

			// Check if option matches any command option
			if opt := command.FindOptionByLongName(cur); opt != nil {
				boolean := true
				if opt.Field().Type == reflect.TypeOf(boolean) {
					continue // If option is boolean, don't skip an argument
				}
				i++ // Else skip next arg in input
				continue
			}
		}

		// Safety check
		if input[i] == "" || input[i] == " " {
			continue
		}

		remain = append(remain, input[i])
	}

	return
}

// MenuHint - Returns the Hint for a given menu context
func MenuHint(args []string, current []rune) (hint []rune) {
	// Menu hints
	var mainMenuHint = tui.DIM + "Client console. Type 'help' or 'help <category>' for more.)"
	var ghostMenuHint = tui.DIM + fmt.Sprintf("In Ghost implant: (adddress: %s%s%s%s). Type 'help' for more.",
		tui.BLUE, commands.Context.Ghost.RemoteAddress, tui.RESET, tui.DIM)

	switch *commands.Context.Menu {
	case commands.MAIN_CONTEXT, commands.MODULE_CONTEXT:
		hint = []rune(menuCatHint + mainMenuHint)
		return
	case commands.GHOST_CONTEXT:
		hint = []rune(menuCatHint + ghostMenuHint)
		return
	}
	return current
}

func handleSpecialCommands(args []string, current []rune) (hint []rune) {
	if args[0] == "exit" {
		hint = []rune(commandCatHint + "Exit the Wiregost console")
		return
	}

	// If goes here, return hint as it was when passed
	return current
}

// envVarHint - Yields hints for environment variables
func envVarHint(args []string, last []rune) (hint []rune) {

	// Trim last in case its a path with multiple vars
	allVars := strings.Split(string(last), "/")
	lastVar := allVars[len(allVars)-1]

	// Base hint
	hint = []rune(envCatHint + lastVar)

	envVar := strings.TrimPrefix(lastVar, "$")

	if v, ok := util.SystemEnv[envVar]; ok {
		if v != "" {
			hintStr := string(hint) + " => " + util.SystemEnv[envVar]
			hint = []rune(hintStr)
		}
	}
	return
}
