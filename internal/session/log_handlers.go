package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

func (s *Session) setLogLevel(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: use 'log level <level>'\n", tui.RED, tui.RESET)
		return
	}
	s.send(cmd)
	log := <-s.logReqs
	fmt.Println()
	fmt.Println(log.Log)
}

func (s *Session) logShow(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: use 'log show <component>'\n", tui.RED, tui.RESET)
		return
	}
	s.send(cmd)
	logs := <-s.logReqs
	for _, l := range logs.Logs {
		fmt.Println(l)
	}
}
