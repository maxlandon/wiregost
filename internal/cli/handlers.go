package cli

import (
	"fmt"
	"strings"
)

// MODULE HANDLERS
//---------------------------------------------------------------------------

func (s *Session) UseModule(cmd []string) {
	s.Send(cmd)
	mod := <-moduleReqs
	s.moduleContext = mod.ModuleName
	CurrentModule = mod.ModuleName
	// Add code to change current module in the prompt
}

func (s *Session) GetModuleOptions() {
	testOptions := strings.Fields("show options")
	s.Send(testOptions)

	mod := <-moduleReqs
	fmt.Println(mod.Options)
	// Make a ShowOptions() func (s *Session)tion here
}

func (s *Session) GetModuleList(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs

	list := mod.Modules
	fmt.Println(list)
}

func (s *Session) SetModuleOption(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs
	fmt.Println(mod)
	// Add some verification that option is correctly set here.
}

func (s *Session) SetAgent(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs
	fmt.Println(mod)
	// Add some verification that agent is correctly set here.
}

func (s *Session) RunModule(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs
	fmt.Println(mod)
	// Add some verification that agent is correctly set here.
}

// AGENT HANDLERS
//---------------------------------------------------------------------------
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

// LOG HANDLERS
//---------------------------------------------------------------------------

func (s *Session) LogLevel(cmd []string) {
	// Send(cmd)
	log := <-logReqs
	fmt.Println(log)
	// Handle change of state here
}

func (s *Session) LogShow(cmd []string) {
	// Send(cmd)
	log := <-logReqs
	fmt.Println(log)
	// Handle printing the logs here
}

// Handle all log messages coming from the server
func (s *Session) LogListen() {
	go func() {
		for {
			msg := <-logReqs
			fmt.Println(msg)
		}
	}()
}

// WORKSPACE HANDLERS
//---------------------------------------------------------------------------

func (s *Session) WorkspaceList(cmd []string) {
	s.Send(cmd)
	workspace := <-workspaceReqs
	fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceSwitch(cmd []string) {
	s.currentWorkspace = cmd[2]
	CurrentWorkspace = cmd[2]
	// fmt.Println(CurrentWorkspace)
	s.Send(cmd)
	workspace := <-workspaceReqs
	s.CurrentWorkspaceId = workspace.WorkspaceId
	// fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceDelete(cmd []string) {
	// Send(cmd)
	workspace := <-workspaceReqs
	fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceNew(cmd []string) {
	s.Send(cmd)
	// workspace := <-workspaceReqs
	// fmt.Println(workspace)
	// Handle change of state here
}

// STACK HANDLERS
//---------------------------------------------------------------------------

func (s *Session) StackShow(cmd []string) {
	// Send(cmd)
	stack := <-stackReqs
	fmt.Println(stack)
	// handle change of state here
}

func (s *Session) StackPop(cmd []string) {
	// Send(cmd)
	stack := <-stackReqs
	fmt.Println(stack)
	// handle change of state here
}

// SERVER HANDLERS
//---------------------------------------------------------------------------

func (s *Session) ServerConnect() { // This command will be used to send connection demands to a server.
	// Send(cmd)
	server := <-serverReqs
	fmt.Println(server)
}
