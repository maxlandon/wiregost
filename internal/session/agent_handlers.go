package session

import "fmt"

func (s *Session) ListAgents(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion for listing agents
}

func (s *Session) InfoAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion to print info.
}

func (s *Session) RemoveAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) DownloadAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) UploadAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) ExecuteShellCodeAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) KillAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) SetAgentOption(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) GetAgentShell(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) BackMainMenu(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Handle change of state here
}

func (s *Session) MainMenu(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Handle change of state here
}
