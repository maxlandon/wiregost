package session

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

// Prompt object
type Prompt struct {
	// Prompt strings
	PromptVariable  string
	DefaultPrompt   string
	ModulePrompt    string
	AgentPrompt     string
	CompilerPrompt  string
	MultilinePrompt string
	// Prompt variables
	CurrentWorkspace *string
	CurrentModule    *string
	MenuContext      *string
	// Other prompt variables
	serverIP *string
	// Callbacks and colors
	PromptCallbacks map[string]func() string
	effects         map[string]string
}

func newPrompt(s *Session) Prompt {
	// These are here because if colors are disabled, we need the updated tui.* variable
	prompt := Prompt{
		// Prompt strings
		PromptVariable:  "$",
		DefaultPrompt:   "{bddg}{fw}@{lb}{localip}{fw} {reset} {dim}in {b}{workspace} {server}",
		ModulePrompt:    "{bddg}{fw}@{lb}{localip}{fw} {reset} {dim}in {b}{workspace} {server} {fw}=>{reset} post({y}{mod}{reset})",
		AgentPrompt:     "{bddg}{fw}@{lb}{localip}{fw} {reset} {dim}in {b}{workspace} {server} {fw}=>{reset} agent[{bold}{db}{agent}{reset}]",
		CompilerPrompt:  "{bddg}{fw}@{lb}{localip}{fw} {reset} {dim}in {b}{workspace} {server} {fw}=>{reset} [{bold}{y}Compiler{reset}]",
		MultilinePrompt: "> ",
		// Prompt variabes
		CurrentWorkspace: &s.currentWorkspace,
		CurrentModule:    &s.currentModule,
		MenuContext:      &s.menuContext,
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
	prompt.PromptCallbacks = map[string]func() string{
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
			return *prompt.CurrentWorkspace
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
		// CurrentModule
		"{mod}": func() string {
			return *prompt.CurrentModule
		},
		// Current agent
		"{agent}": func() string {
			return s.currentAgentID.String()
		},
		// Server state
		"{server}": func() string {
			if s.serverRunning {
				return fmt.Sprintf("%s(%son%s)", tui.RESET, tui.GREEN, tui.RESET)
			} else {
				return fmt.Sprintf("%s(%soff%s)", tui.RESET, tui.RED, tui.RESET)
			}
		},
	}

	return prompt
}

func (p Prompt) render() (first string, multi string) {

	var prompt string

	// Current module does not depend on context...
	if *p.CurrentModule != "" {
		prompt = p.ModulePrompt
	} else {
		prompt = p.DefaultPrompt
	}
	// ... and is overidden by the context string if needed.
	if *p.MenuContext == "compiler" {
		prompt = p.CompilerPrompt
	}
	// ... or overriden by the context agent if needed.
	if *p.MenuContext == "agent" {
		prompt = p.AgentPrompt
	}

	multiline := p.MultilinePrompt

	for tok, effect := range p.effects {
		prompt = strings.Replace(prompt, tok, effect, -1)
		multiline = strings.Replace(multiline, tok, effect, -1)
	}

	for tok, cb := range p.PromptCallbacks {
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
	p, _ := prompt.render()
	_, m := prompt.render()
	fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
}
