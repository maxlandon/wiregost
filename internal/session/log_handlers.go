package session

import (
	"fmt"
)

func (s *Session) setLogLevel(cmd []string) {
	s.send(cmd)
	log := <-s.logReqs
	fmt.Println()
	fmt.Println(log.Log)
}

func (s *Session) logShow(cmd []string) {
	s.send(cmd)
	logs := <-s.logReqs
	for _, l := range logs.Logs {

		fmt.Println(l)
	}
}
