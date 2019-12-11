package session

// This file contains all the code necessary for the shell prompt.
// This includes environment variabes (real-time computed), constructors,
// and rendering functions.

import (
	"net"
	"os"
	"strings"

	"github.com/maxlandon/wiregost/internal/session/core"

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
	PromptCallbacks = map[string]func(s *Session) string{
		"{iface}": func(s *Session) string {
			return "192.168.1.0/24"
		},
		// Working directory
		"{pwd}": func(s *Session) string {
			pwd, _ := os.Getwd()
			return pwd
		},
		// Current Workspace
		"{workspace}": func(s *Session) string {
			return "Fixed_Workspace"
		},
		// Local IP address
		"{localip}": func(s *Session) string {
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
		"{blink}": core.BLINK,
	}
	return Prompt{}
}

func (p Prompt) Render(s *Session) (first string, multi string) {
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
		prompt = strings.Replace(prompt, tok, cb(s), -1)
		multiline = strings.Replace(multiline, tok, cb(s), -1)
	}

	// make sure an user error does not screw all terminal
	if !strings.HasPrefix(prompt, tui.RESET) {
		prompt += tui.RESET
	}

	return prompt, multiline
}
