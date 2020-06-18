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

package console

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/lmorg/readline"

	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/context"
)

var (
	// Prompt - The prompt object used by the console
	Prompt *prompt
)

// prompt - Stores all variables necessary to the console prompt
type prompt struct {
	// Strings
	BaseMain       string
	BaseModule     string
	BaseGhost      string
	CustomMain     string
	CustomGhost    string
	MultilineEmacs string
	MultilineVim   string
	// Callbacks & Colors
	Callbacks map[string]func() string
	Effects   map[string]string
}

// SetPrompt - Initializes the Prompt object
func (c *console) SetPrompt() {

	// Initialize
	Prompt = &prompt{
		BaseMain:       "{bddg}{fw}@{lb}{serverip} {reset} {dim}in {workspace} {reset}({g}{listeners}{fw},{r}{agents}{fw})",
		BaseModule:     " {dim}=>{reset} {type}({mod})",
		BaseGhost:      "{bddg}{fw}agent[{lb}{agent}]{reset} ",
		CustomMain:     "",
		CustomGhost:    "{dim}as {user}{bold}{y}@{reset}{host}/{rpwd} {dim}in {workspace}",
		MultilineVim:   "{vim} > ",
		MultilineEmacs: " > ",
	}

	setCallbacks(Prompt)

	return
}

// setCallbacks - Initializes all callbacks for prompt
func setCallbacks(prompt *prompt) {

	// Colors
	prompt.Effects = map[string]string{
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
	prompt.Callbacks = map[string]func() string{
		// Vim mode
		"{vim}": func() string {
			return ""
		},

		// Working directory
		"{pwd}": func() string {
			pwd, _ := os.Getwd()
			return pwd
		},
		// Current Workspace
		"{workspace}": func() string {
			return tui.Blue(context.Context.Workspace.Name)
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
			return assets.ServerConfig.LHost
		},
		// Listeners
		"{listeners}": func() string {
			listeners := strconv.Itoa(context.Context.Jobs)
			return listeners
		},
		// Agents
		"{agents}": func() string {
			agents := strconv.Itoa(context.Context.Ghosts)
			return agents
		},
		// Current Module type
		"{type}": func() string {
			if len(context.Context.Module.Path) != 0 {
				switch strings.Split(context.Context.Module.Path, "/")[0] {
				case "post":
					return "post"
				case "exploit":
					return "exploit"
				case "auxiliary":
					return "auxiliary"
				case "payload":
					return "payload"
				}
			}
			return ""
		},
		// CurrentModule
		"{mod}": func() string {
			if len(context.Context.Module.Path) != 0 {
				return tui.Red(tui.Bold(context.Context.Module.Path)) + tui.RESET
			}
			return ""
			// return tui.Yellow(*prompt.currentModule) + tui.RESET
		},
		// Current agent
		"{agent}": func() string {
			return context.Context.Ghost.Name
		},
		// Agent username
		"{user}": func() string {
			return tui.RESET + tui.Bold(context.Context.Ghost.Username)
		},
		// Agent hostname
		"{host}": func() string {
			return tui.Bold(context.Context.Ghost.Hostname) + tui.RESET
		},
		// Agent cwd
		"{rpwd}": func() string {
			return tui.Blue(context.Context.Ghost.Pwd) + tui.RESET
		},
		// agent user ID
		"{uid}": func() string {
			return context.Context.Ghost.UID
		},
		// agent user group ID
		"{gid}": func() string {
			return context.Context.Ghost.GID
		},
		// agent process ID
		"{pid}": func() string {
			return strconv.Itoa(int(context.Context.Ghost.PID))
		},
		// agent C2 protocol
		"{transport}": func() string {
			// return context.Context.Ghost.Transports[0].LHost + ":" + strconv.Itoa(int(context.Context.Ghost.Transports[0].LPort))
			return ""
		},
		// agent remote host:port address
		"{address}": func() string {
			// return context.Context.Ghost.Transports[0].RHost + ":" + strconv.Itoa(int(context.Context.Ghost.Transports[0].RPort))
			return ""
		},
		// agent target OS
		"{os}": func() string {
			return context.Context.Ghost.OS
		},
		// agent target CPU Arch
		"{arch}": func() string {
			return context.Context.Ghost.Arch
		},
	}

}

// render - Computes all variables and outputs prompt
func (p *prompt) render() (prompt string, multi string) {

	ctx := context.Context

	switch p.CustomMain {
	// No custom prompt provided, use base
	case "":
		if len(ctx.Module.Path) != 0 {
			// if len(commands.Context.Module.Path) != 0 {
			prompt = p.BaseMain + p.BaseModule
		} else {
			prompt = p.BaseMain
		}
		if ctx.Menu == context.GHOST_CONTEXT {
			// Check custom implant prompt provided
			if p.CustomGhost == "" {
				prompt = p.BaseGhost
			} else {
				prompt = p.CustomGhost
			}
		}

	// Custom provided, use it
	default:
		if len(ctx.Module.Path) != 0 {
			// if len(commands.Context.Module.Path) != 0 {
			prompt = p.CustomMain + p.BaseModule
		} else {
			prompt = p.CustomMain
		}
		if ctx.Menu == context.GHOST_CONTEXT {
			// Check custom implant prompt provided
			if p.CustomGhost == "" {
				prompt = p.BaseGhost + p.CustomGhost // We keep agent[NAME] in both cases
			} else {
				prompt = p.BaseGhost + p.CustomGhost
			}
		}
	}

	// Set multiline based on input mode
	multiline := p.MultilineEmacs

	for tok, effect := range p.Effects {
		prompt = strings.Replace(prompt, tok, effect, -1)
		multiline = strings.Replace(multiline, tok, effect, -1)
	}

	for tok, cb := range p.Callbacks {
		prompt = strings.Replace(prompt, tok, cb(), -1)
		multiline = strings.Replace(multiline, tok, cb(), -1)
	}

	// make sure an user error does not screw all terminal
	if !strings.HasPrefix(prompt, tui.RESET) {
		prompt += tui.RESET
	}

	return prompt, multiline
}

// RefreshPrompt - Recompute prompt
func RefreshPrompt(prompt *prompt, input *readline.Instance) {
	p, _ := prompt.render()
	_, m := prompt.render()
	// fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
}
