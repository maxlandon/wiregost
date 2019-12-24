package cli

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

const (
	PromptVariable  = "$"
	DefaultPrompt   = "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > {b}{pwd} {reset}"
	ModulePrompt    = "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > {b}{pwd} {reset}post({r}{bold}{mod}{reset})"
	MultilinePrompt = "{g}> {reset}"
)

// Current shell state variables
var CurrentWorkspace = "default"
var CurrentModule string
var serverIp string

// Prompt real-time Environment variables
var (
	effects         = map[string]string{}
	PromptCallbacks = map[string]func() string{
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
			return CurrentWorkspace
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
			return CurrentModule
		},
	}
)

// Prompt object
type Prompt struct {
}

func NewPrompt() Prompt {
	// these are here because if colors are disabled,
	// we need the updated tui.* variable
	effects = map[string]string{
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
	return Prompt{}
}

func (p Prompt) Render() (first string, multi string) {

	var prompt string

	// Set prompt depending on context
	if CurrentModule != "" {
		prompt = ModulePrompt
	} else {
		prompt = DefaultPrompt
	}
	multiline := MultilinePrompt

	for tok, effect := range effects {
		prompt = strings.Replace(prompt, tok, effect, -1)
		multiline = strings.Replace(multiline, tok, effect, -1)
	}

	for tok, cb := range PromptCallbacks {
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
func Refresh(prompt Prompt, input *readline.Instance) {
	p, _ := prompt.Render()
	_, m := prompt.Render()
	// p, _ := s.parseEnvTokens(s.Prompt.Render(s))
	fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
	input.Refresh()
}
