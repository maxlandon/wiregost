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
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/maxlandon/wiregost/client/commands"
)

// These functions are just shorthands for checking various conditions on the input line. They make the main function more readable, which might be
// useful, should a logic error pop somewhere.

// [ Menus ]
// Is the input line is either empty, or without any detected command ?
func noCommandOrEmpty(args []string, last []rune, command *flags.Command) bool {
	if len(args) == 0 || len(args) == 1 && command == nil {
		return true
	}
	return false
}

// [ Commands ]
// detectedCommand - Returns the base command from parser if detected, depending on context
func detectedCommand(args []string) *flags.Command {

	cmds := commands.CommandsByContext()
	var command *flags.Command
	for _, cmd := range cmds {
		if cmd.Name == args[0] {
			command = cmd
		}
	}
	return command
}

// is the command a special command, usually not handled by parser ?
func isSpecialCommand(args []string, command *flags.Command) bool {

	// If command is not nil, return
	if command == nil {
		// Shell
		if args[0] == "!" {
			return true
		}
		// Exit
		if args[0] == "exit" {
			return true
		}
		return false
	}
	return false
}

// The commmand has been found
func commandFound(command *flags.Command) bool {
	if command != nil {
		return true
	}
	return false
}

// [ SubCommands ]
// Does the command have subcommands ?
func hasSubCommands(command *flags.Command, args []string) bool {
	if len(args) > 2 || command == nil {
		return false
	}

	if len(command.Commands()) != 0 {
		return true
	}

	return false
}

// Does the input has a subcommand in it ?
func subCommandFound(last []rune, args []string, command *flags.Command) (sub *flags.Command, ok bool) {
	if len(args) < 2 || command == nil {
		return nil, false
	}

	sub = command.Find(args[1])
	if sub != nil {
		return sub, true
	}

	return nil, false
}

// Is the last input PRECISELY a subcommand. This is used as a brief hint for the subcommand
func lastIsSubCommand(last []rune, command *flags.Command) bool {
	if sub := command.Find(string(last)); sub != nil {
		return true
	}
	return false
}

// [ Arguments ]
// Does the command have arguments ?
func hasArgs(command *flags.Command) bool {
	if len(command.Args()) != 0 {
		return true
	}
	return false
}

// [ Options ]
// optionsAsked - Does the user asks for options ?
func optionsAsked(args []string, last []rune, command *flags.Command) bool {
	if len(args) > 2 && (strings.HasPrefix(string(last), "-") || strings.HasPrefix(string(last), "--")) {
		return true
	}
	return false
}

// Is the last input argument is a dash ?
func isOptionDash(args []string, last []rune) bool {
	if len(args) > 2 && (strings.HasPrefix(string(last), "-") || strings.HasPrefix(string(last), "--")) {
		return true
	}
	return false
}

// optionIsAlreadySet - Detects in input if an option is already set
func optionIsAlreadySet(args []string, last []rune, opt *flags.Option) bool {
	return false
}

// Check if option type allows for repetition
func optionNotRepeatable(opt *flags.Option) bool {
	return true
}

// [ Option Values ]
// Is the last input word an option name (--option) ?
func optionArgRequired(args []string, last []rune, group *flags.Group) (opt *flags.Option, yes bool) {

	var lastItem string
	var lastOption string
	var option *flags.Option

	// Check for last two arguments in input
	if args[len(args)-1] == "" {
		lastItem = args[len(args)-1]

		if strings.HasPrefix(args[len(args)-2], "-") || strings.HasPrefix(args[len(args)-2], "--") {
			lastOption = strings.TrimPrefix(args[len(args)-2], "--")
			lastOption = strings.TrimPrefix(lastOption, "-")

			if opt := group.FindOptionByLongName(lastOption); opt != nil {
				option = opt
			}
		}
	} else {
		lastItem = args[len(args)-1]
		if strings.HasPrefix(args[len(args)-1], "-") || strings.HasPrefix(args[len(args)-1], "--") {
			lastOption = strings.TrimPrefix(args[len(args)-1], "--")
			lastOption = strings.TrimPrefix(lastOption, "-")

			if opt := group.FindOptionByLongName(lastOption); opt != nil {
				option = opt
			}
		}
	}

	// If option is found, and we still are in writing the argument
	if (lastItem == "" && option != nil) || option != nil {
		// Check if option is a boolean, if yes return false
		boolean := true
		if option.Field().Type == reflect.TypeOf(boolean) {
			return nil, false
		}

		// Check this recursion and its effects !!!!!
		if len(group.Groups()) != 0 {
			for _, grp := range group.Groups() {
				opt, found := optionArgRequired(args, last, grp)
				if found {
					return opt, found
				}
			}
			return nil, false
		}
		return option, true
	}

	// Check for previous argument
	if lastItem != "" && option == nil {
		if strings.HasPrefix(args[len(args)-2], "-") || strings.HasPrefix(args[len(args)-2], "--") {

			lastOption = strings.TrimPrefix(args[len(args)-2], "--")
			lastOption = strings.TrimPrefix(lastOption, "-")

			if opt := group.FindOptionByLongName(lastOption); opt != nil {
				option = opt
				return option, true
			}

		}
	}

	return nil, false
}

// [ Other ]
// Does the user asks for Environment variables ?
func envVarAsked(args []string, last []rune) bool {

	// Check if the current word is an environment variable, or if the last part of it is a variable
	if len(last) > 1 && strings.HasPrefix(string(last), "$") {
		if strings.LastIndex(string(last), "/") < strings.LastIndex(string(last), "$") {
			return true
		}
		return false
	}

	// Check if env var is asked in a path or something
	if len(last) > 1 {
		// If last is a path, it cannot be an env var anymore
		if last[len(last)-1] == '/' {
			return false
		}

		if last[len(last)-1] == '$' {
			return true
		}
	}

	// If we are at the beginning of an env var
	if len(last) > 0 && last[len(last)-1] == '$' {
		return true
	}

	return false
}
