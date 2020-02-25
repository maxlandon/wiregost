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
	"strconv"
	"strings"

	// 3rd party
	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

// Prompt object
type Prompt struct {
	// Prompt strings
	baseMain       string
	baseAgent      string
	module         string
	agent          string
	customMain     string
	customAgent    string
	multilineVim   string
	multilineEmacs string
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

func newPrompt(c *Console, customMain string, customAgent string) Prompt {
	// These are here because if colors are disabled, we need the updated tui.* variable
	prompt := Prompt{
		// Prompt strings
		baseMain:       "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {reset}({g}{listeners}{fw},{r}{agents}{fw})",
		module:         " =>{reset} {type}({mod})",
		baseAgent:      "{dim}as {user}{bold}{y}@{reset}{host}/{rpwd} {dim}in {workspace}",
		agent:          "{bddg}{fw}agent[{lb}{agent}]{reset} ",
		customMain:     customMain,
		customAgent:    customAgent,
		multilineVim:   "{vim} > ",
		multilineEmacs: " > ",
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
		"{ly}":    "\033[38;5;187m",
	}
	// Callbacks
	prompt.promptCallbacks = map[string]func() string{
		// Vim mode
		"{vim}": func() string {
			switch c.vimMode {
			case "insert":
				return tui.Yellow("[I]")
			case "normal":
				return tui.Blue("[N]")
			}
			return ""
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
			return c.server.Config.LHost
		},
		// Listeners
		"{listeners}": func() string {
			listeners := strconv.Itoa(c.listeners)
			return listeners
		},
		// Agents
		"{agents}": func() string {
			agents := strconv.Itoa(c.ghosts)
			return agents
		},
		// Current Module type
		"{type}": func() string {
			switch strings.Split(c.currentModule, "/")[0] {
			case "post":
				return "post"
			case "exploit":
				return "exploit"
			case "auxiliary":
				return "auxiliary"
			case "payload":
				return "payload"
			}
			return ""
		},
		// CurrentModule
		"{mod}": func() string {
			mod := strings.Join(strings.Split(*prompt.currentModule, "/")[1:], "/")
			return tui.Red(tui.Bold(mod)) + tui.RESET
			// return tui.Yellow(*prompt.currentModule) + tui.RESET
		},
		// Current agent
		"{agent}": func() string {
			return c.CurrentAgent.Name
		},
		// Agent username
		"{user}": func() string {
			return tui.RESET + tui.Bold(c.CurrentAgent.Username)
		},
		// Agent hostname
		"{host}": func() string {
			return tui.Bold(c.CurrentAgent.Hostname) + tui.RESET
		},
		// Agent cwd
		"{rpwd}": func() string {
			return tui.Blue(c.AgentPwd) + tui.RESET
		},
		// agent user ID
		"{uid}": func() string {
			return c.CurrentAgent.UID
		},
		// agent user group ID
		"{gid}": func() string {
			return c.CurrentAgent.GID
		},
		// agent process ID
		"{pid}": func() string {
			return strconv.Itoa(int(c.CurrentAgent.PID))
		},
		// agent C2 protocol
		"{transport}": func() string {
			return c.CurrentAgent.Transport
		},
		// agent remote host:port address
		"{address}": func() string {
			return c.CurrentAgent.RemoteAddress
		},
		// agent target OS
		"{os}": func() string {
			return c.CurrentAgent.OS
		},
		// agent target CPU Arch
		"{arch}": func() string {
			return c.CurrentAgent.Arch
		},
	}

	return prompt
}

func (p Prompt) render(vimMode bool) (first string, multi string) {

	var prompt string

	switch p.customMain {
	// No custom prompt provided, use base
	case "":
		if *p.currentModule != "" {
			prompt = p.baseMain + p.module
		} else {
			prompt = p.baseMain
		}
		if *p.menu == "agent" {
			// Check custom implant prompt provided
			if p.customAgent == "" {
				prompt = p.agent + p.baseAgent
			} else {
				prompt = p.agent + p.customAgent
			}
		}

	// Custom provided, use it
	default:
		if *p.currentModule != "" {
			prompt = p.customMain + p.module
		} else {
			prompt = p.customMain
		}
		if *p.menu == "agent" {
			// Check custom implant prompt provided
			if p.customAgent == "" {
				prompt = p.agent + p.baseAgent
			} else {
				prompt = p.agent + p.customAgent
			}
		}
	}

	// Set multiline based on input mode
	multiline := p.multilineEmacs
	if vimMode {
		multiline = p.multilineVim
	}

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
	return prompt, multiline
}

// Refresh prompt
func refreshPrompt(prompt Prompt, input *readline.Instance) {
	p, _ := prompt.render(input.IsVimMode())
	_, m := prompt.render(input.IsVimMode())
	fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
}
