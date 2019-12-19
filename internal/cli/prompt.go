package cli

import (
	"net"
	"os"
	"strings"

	"github.com/evilsocket/islazy/tui"
)

const (
	PromptVariable  = "$"
	DefaultPrompt   = "{bdg}{y}{localip} {fb}|{fw} {workspace} {reset} > {b}{pwd} {reset}"
	MultilinePrompt = "{g}> {reset}"
)

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
			return "Fixed_Workspace"
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
	// found, prompt := s.Env.Get(PromptVariable)		// Used if Custom prompt is saved in Env Config
	// if !found {
	//     prompt = DefaultPrompt
	// }
	//
	prompt := DefaultPrompt
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
