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
	CompilerPrompt  string
	MultilinePrompt string
	// Prompt variables
	CurrentWorkspace *string
	CurrentModule    *string
	MenuContext      *string
	// Other prompt variables
	serverIp *string
	// Callbacks and colors
	PromptCallbacks map[string]func() string
	effects         map[string]string
}

func NewPrompt(s *Session) Prompt {
	// these are here because if colors are disabled,
	// we need the updated tui.* variable
	prompt := Prompt{
		// Prompt strings
		PromptVariable:  "$",
		DefaultPrompt:   "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > {dim}in {b}{pwd} {reset}",
		ModulePrompt:    "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > {reset}post({r}{bold}{mod}{reset}) {dim}in {b}{pwd} ",
		CompilerPrompt:  "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > [{bold}{y}Compiler{reset}] {dim}in {b}{pwd} {reset} ",
		MultilinePrompt: "{g}> {reset}",
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
				networkIp, ok := addr.(*net.IPNet)
				if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
					ip = networkIp.IP.String()
				}
			}
			return ip
		},
		// CurrentModule
		"{mod}": func() string {
			return *prompt.CurrentModule
		},
	}

	return prompt
}

func (p Prompt) Render() (first string, multi string) {

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
func RefreshPrompt(prompt Prompt, input *readline.Instance) {
	p, _ := prompt.Render()
	_, m := prompt.Render()
	fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
}
