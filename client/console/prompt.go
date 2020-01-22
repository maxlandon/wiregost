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

package console

import (
	// Standard
	"fmt"
	"net"
	"os"
	"strings"

	// 3rd party
	"github.com/evilsocket/islazy/tui"
)

// Prompt object
type Prompt struct {
	// Prompt strings
	base      string
	module    string
	agent     string
	compiler  string
	multiline string
	// Prompt variables
	workspace     *string
	currentModule *string
	menu          *string
	// Other prompt variables
	serverIP *string
	// Callbacks and colors
	promptCallbacks map[string]func() string
	effects         map[string]string
}

func newPrompt(c *Console) Prompt {
	// These are here because if colors are disabled, we need the updated tui.* variable
	prompt := Prompt{
		// Prompt strings
		base:      "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {server} > ",
		module:    "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {server} {fw}=>{reset} post({mod})",
		agent:     "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {server} {fw}=>{reset} agent[{db}{agent}]",
		compiler:  "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {server} {fw}=>{reset} [{bold}{y}Compiler{reset}]",
		multiline: "> ",
		// Prompt variabes
		workspace:     &c.currentWorkspace.Name,
		currentModule: &c.currentModule,
		menu:          &c.menuContext,
	}
	// Colors
	prompt.effects = map[string]string{
		"{bold}":  tui.BOLD,
		"{dim}":   tui.DIM,
		"{r}":     tui.RED,
		"{g}":     tui.GREEN,
		"{b}":     tui.BLUE,
		"{y}":     tui.YELLOW,
		"{fb}":    tui.FOREBLACK,
		"{fw}":    tui.FOREWHITE,
		"{bdg}":   tui.BACKDARKGRAY,
		"{br}":    tui.BACKRED,
		"{bg}":    tui.BACKGREEN,
		"{by}":    tui.BACKYELLOW,
		"{blb}":   tui.BACKLIGHTBLUE,
		"{reset}": tui.RESET,

		// Custom colors:
		"{blink}": "\033[5m",
		"{lb}":    "\033[38;5;117m",
		"{db}":    "\033[38;5;24m",
		"{bddg}":  "\033[48;5;237m",
	}
	// Callbacks
	prompt.promptCallbacks = map[string]func() string{
		"{iface}": func() string {
			return "192.168.1.0/24"
		},
		// Working directory
		"{pwd}": func() string {
			pwd, _ := os.Getwd()
			return pwd
		},
		// Current Workspace
		"{workspace}": func() string {
			return tui.Blue(*prompt.workspace)
		},
		// Local IP address
		"{localip}": func() string {
			addrs, _ := net.InterfaceAddrs()
			var ip string
			for _, addr := range addrs {
				network, ok := addr.(*net.IPNet)
				if ok && !network.IP.IsLoopback() && network.IP.To4() != nil {
					ip = network.IP.String()
				}
			}
			return ip
		},
		"{serverip}": func() string {
			return c.endpointPublicIP
		},
		// CurrentModule
		"{mod}": func() string {
			return tui.Yellow(*prompt.currentModule) + tui.RESET
		},
		// Current agent
		"{agent}": func() string {
			return tui.Bold(c.currentAgentID.String()) + tui.RESET
		},
		// Server state
		"{server}": func() string {
			if c.serverRunning {
				return fmt.Sprintf("%s(%son%s)", tui.RESET, tui.GREEN, tui.RESET)
			} else {
				return fmt.Sprintf("%s(%soff%s)", tui.RESET, tui.RED, tui.RESET)
			}
		},
	}

	return prompt
}

func (p Prompt) render() (first string) {

	var prompt string

	// Current module does not depend on context...
	if *p.currentModule != "" {
		prompt = p.module
	} else {
		prompt = p.base
	}
	// ... and is overidden by the context string if needed.
	if *p.menu == "compiler" {
		prompt = p.compiler
	}
	// ... or overriden by the context agent if needed.
	if *p.menu == "agent" {
		prompt = p.agent
	}

	multiline := p.multiline

	for tok, effect := range p.effects {
		prompt = strings.Replace(prompt, tok, effect, -1)
		multiline = strings.Replace(multiline, tok, effect, -1)
	}

	for tok, cb := range p.promptCallbacks {
		prompt = strings.Replace(prompt, tok, cb(), -1)
		multiline = strings.Replace(multiline, tok, cb(), -1)
	}

	// make sure an user error does not screw all terminal
	if !strings.HasPrefix(prompt, tui.RESET) {
		prompt += tui.RESET
	}
	return prompt
}

// func (p Prompt) render() (first string, multi string) {
//
//         var prompt string
//
//         // Current module does not depend on context...
//         if *p.currentModule != "" {
//                 prompt = p.module
//         } else {
//                 prompt = p.base
//         }
//         // ... and is overidden by the context string if needed.
//         if *p.menu == "compiler" {
//                 prompt = p.compiler
//         }
//         // ... or overriden by the context agent if needed.
//         if *p.menu == "agent" {
//                 prompt = p.agent
//         }
//
//         multiline := p.multiline
//
//         for tok, effect := range p.effects {
//                 prompt = strings.Replace(prompt, tok, effect, -1)
//                 multiline = strings.Replace(multiline, tok, effect, -1)
//         }
//
//         for tok, cb := range p.promptCallbacks {
//                 prompt = strings.Replace(prompt, tok, cb(), -1)
//                 multiline = strings.Replace(multiline, tok, cb(), -1)
//         }
//
//         // make sure an user error does not screw all terminal
//         if !strings.HasPrefix(prompt, tui.RESET) {
//                 prompt += tui.RESET
//         }
//         return prompt, multiline
// }

// Refresh prompt
// func refreshPrompt(prompt Prompt, input *readline.Instance) {
//         p, _ := prompt.render()
//         _, m := prompt.render()
//         fmt.Println()
//         fmt.Println(p)
//         input.SetPrompt(m)
// }
