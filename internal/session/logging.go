package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
)

// Log needs to be exported because of log package conflict.
func Log(event messages.LogEvent) {
	// Set colors and signs for levels
	switch event.Level {
	case "debug":
		fmt.Printf("[%sdebug%s] %s \n", tui.DIM, tui.RESET, event.Message)
	case "info":
		fmt.Printf("%s[-]%s %s \n", tui.BOLD, tui.RESET, event.Message)
	case "warning":
		fmt.Printf("%s[!]%s %s \n", tui.YELLOW, tui.RESET, event.Message)
	case "error":
		fmt.Printf("%s%s[!]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event.Message)
	case "fatal":
		fmt.Printf("%s%s[FATAL]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event.Message)
	case "panic":
		fmt.Printf("%s%s[PANIC]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event.Message)
	}
}
