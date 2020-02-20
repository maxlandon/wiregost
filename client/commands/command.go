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
	"time"

	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

var (
	defaultTimeout   = 30 * time.Second
	stdinReadTimeout = 10 * time.Millisecond
)

// RPCServer - Function used to send/recv envelopes
type RPCServer func(*ghostpb.Envelope, time.Duration) chan *ghostpb.Envelope

// Command is a set of commands dedicated to a single function (workspace, agents, etc)
type Command struct {
	Name        string
	Help        string
	SubCommands []string
	Args        []*CommandArg
	Handle      func(*Request) error
}

// CommandArg is an argument to a command/subcommand, like host-id when searching for hosts
type CommandArg struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Length      int
}

// commandMap maps commands to an appropriate menu context
var commandList = map[string]*Command{}
var commandMap = map[string]map[string]*Command{}

// Request creates a request from a command, passing it all necessary shell context
type Request struct {
	// Command
	Command *Command
	Args    []string

	// Shell context
	context *ShellContext
}

// NewRequest creates a request, with the shell context
func NewRequest(cmd *Command, args []string, shellContext *ShellContext) *Request {
	return &Request{
		Command: cmd,
		Args:    args,
		context: shellContext,
	}
}

// FindCommand finds a commmand for a given menu context
func FindCommand(context, name string) *Command {
	return commandMap[context][name]
}

func AllContextCommands(context string) map[string]*Command {
	return commandMap[context]
}

// AddCommand adds a command/set to a menu context
func AddCommand(context string, cmd *Command) {

	// Check context list exists
	if commandMap == nil {
		commandMap = make(map[string]map[string]*Command)
	}

	// Check map for each context exists
	if commandMap[context] == nil {
		commandMap[context] = map[string]*Command{}
	}

	// Add to context list
	commandMap[context][cmd.Name] = cmd
}
