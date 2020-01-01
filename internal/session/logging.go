package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

func Log(event map[string]string) {
	// Set colors and signs for levels
	switch event["level"] {
	case "debug":
		fmt.Printf("[%sdebug%s] %s \n", tui.DIM, tui.RESET, event["message"])
	case "info":
		fmt.Printf("%s[-]%s %s \n", tui.BOLD, tui.RESET, event["message"])
	case "warning":
		fmt.Printf("%s[!]%s %s \n", tui.YELLOW, tui.RESET, event["message"])
	case "error":
		fmt.Printf("%s%s[!]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event["message"])
	case "fatal":
		fmt.Printf("%s%s[FATAL]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event["message"])
	case "panic":
		fmt.Printf("%s%s[PANIC]%s %s \n", tui.BOLD, tui.RED, tui.RESET, event["message"])
	}
}
