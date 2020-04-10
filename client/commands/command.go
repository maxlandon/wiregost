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

package commands

import (
	"fmt"
	"time"

	"github.com/evilsocket/islazy/tui"
	flags "github.com/jessevdk/go-flags"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

var (
	// CommandMap - keeps a map of all commands available in a given menu context
	CommandMap = map[string][]*flags.Command{}

	// CommandParser - Parses inputs and execute appropriate commands
	MainParser = flags.NewParser(&globalOptions, flags.Default)
	// GhostCommands - Ghost context commands
	GhostParser = flags.NewParser(&globalOptions, flags.Default)

	// Timeouts
	DefaultTimeout   = 30 * time.Second
	StdinReadTimeout = 10 * time.Millisecond
)

// Menu Contexts
const (
	// MAIN_CONTEXT - Available only in main menu
	MAIN_CONTEXT = "main"
	// MODULE_CONTEXT - Available only when a module is loaded
	MODULE_CONTEXT = "module"
	// GHOST_CONTEXT - Available only when interacting with a ghost implant
	GHOST_CONTEXT = "ghost"
)

// GlobalOptions - All options applying to all commands
type GlobalOptions struct{}

var globalOptions GlobalOptions

// RPCServer - Function used to send/recv envelopes
type RPCServer func(*ghostpb.Envelope, time.Duration) chan *ghostpb.Envelope

// CommandsByContext - Returns all commands available to a context
func CommandsByContext() []*flags.Command {
	return CommandMap[*Context.Menu]
}

// OptionByName - Returns an option for a command or a subcommand, identified by name
func OptionByName(context string, command, subCommand, option string) *flags.Option {

	var cmd *flags.Command

	switch context {
	case MAIN_CONTEXT, MODULE_CONTEXT:
		cmd = MainParser.Find(command)
	case GHOST_CONTEXT:
		cmd = GhostParser.Find(command)
	}

	// Base command is found
	if cmd != nil {
		// If options are for a subcommand
		if subCommand != "" && len(cmd.Commands()) != 0 {
			sub := cmd.Find(subCommand)
			if sub != nil {
				for _, opt := range sub.Options() {
					if opt.LongName == option {
						return opt
					}
				}
				return nil
			}
			return nil
		}
		// If subcommand is not asked, return opt for base
		for _, opt := range cmd.Options() {
			if opt.LongName == option {
				return opt
			}
		}
	}
	return nil
}

func ArgumentByName(command *flags.Command, name string) *flags.Arg {
	args := command.Args()
	for _, arg := range args {
		if arg.Name == name {
			return arg
		}
	}

	return nil
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
	// ParserError - Failed to parse some tokens in the input
	ParserError = fmt.Sprintf("%s[Parser Error]%s ", tui.RED, tui.RESET)
	// DBError - Data Service error
	DBError = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)
)
