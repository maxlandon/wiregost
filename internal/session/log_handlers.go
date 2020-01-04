package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

func (s *Session) SetLogLevel(cmd []string) {
	s.Send(cmd)
	log := <-s.logReqs
	fmt.Println()
	fmt.Println(log.Log)
}

func (s *Session) LogShow(cmd []string) {
	s.Send(cmd)
	logs := <-s.logReqs
	for _, l := range logs.Logs {
		switch l["level"] {
		case "debug":
			fmt.Printf("%s%s%s [%sdebug%s] %s \n", tui.DIM, l["time"], tui.RESET, tui.DIM, tui.RESET, l["msg"])
		case "info":
			fmt.Printf("%s%s%s %s[-]%s %s \n", tui.DIM, l["time"], tui.RESET, tui.BOLD, tui.RESET, l["msg"])
		case "warning":
			fmt.Printf("%s%s%s %s[!]%s %s \n", tui.DIM, l["time"], tui.RESET, tui.YELLOW, tui.RESET, l["msg"])
		case "error":
			fmt.Printf("%s%s%s %s%s[!]%s %s \n", tui.DIM, l["time"], tui.RESET, tui.BOLD, tui.RED, tui.RESET, l["msg"])
		case "fatal":
			fmt.Printf("%s%s%s %s%s[FATAL]%s %s \n", tui.DIM, l["time"], tui.RESET, tui.BOLD, tui.RED, tui.RESET, l["msg"])
		case "panic":
			fmt.Printf("%s%s%s %s%s[PANIC]%s %s \n", tui.DIM, l["time"], tui.RESET, tui.BOLD, tui.RED, tui.RESET, l["msg"])
		}
	}
}
