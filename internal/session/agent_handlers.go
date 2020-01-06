package session

import "fmt"

func (s *Session) listAgents(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion for listing agents
}

func (s *Session) infoAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion to print info.
}

func (s *Session) removeAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) downloadAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) uploadAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) executeShellCodeAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) killAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) setAgentOption(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) getAgentShell(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) backMainMenu(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Handle change of state here
}

func (s *Session) mainMenu(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Handle change of state here
}
