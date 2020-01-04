package session

import (
	"bufio"
	"crypto/tls"
	"io"
	"os"
	"strings"
	"time"

	// 3rd party
	"github.com/chzyer/readline"
	"github.com/maxlandon/wiregost/internal/compiler"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/maxlandon/wiregost/internal/modules"
)

type Session struct {
	// Shell
	Shell  *readline.Instance
	prompt Prompt
	// Auth
	user *User
	// Context
	menuContext        string
	currentModule      string
	currentWorkspace   string
	CurrentWorkspaceId int
	// Environmment variables
	Env map[string]string
	// Server connection parameters
	SavedEndpoints  []Endpoint
	CurrentEndpoint Endpoint
	connected       bool
	// Connection
	connection *tls.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	// Response Channels
	moduleReqs    chan modules.ModuleResponse
	agentReqs     chan messages.AgentResponse
	logReqs       chan messages.LogResponse
	workspaceReqs chan messages.WorkspaceResponse
	endpointReqs  chan messages.EndpointResponse
	serverReqs    chan messages.ServerResponse
	stackReqs     chan messages.StackResponse
	compilerReqs  chan compiler.CompilerResponse
	logEventReqs  chan map[string]string
}

func NewSession() *Session {
	session := &Session{
		menuContext: "main",
		Env:         make(map[string]string),
		// Response channels
		moduleReqs:    make(chan modules.ModuleResponse),
		agentReqs:     make(chan messages.AgentResponse),
		logReqs:       make(chan messages.LogResponse),
		workspaceReqs: make(chan messages.WorkspaceResponse),
		endpointReqs:  make(chan messages.EndpointResponse),
		serverReqs:    make(chan messages.ServerResponse),
		stackReqs:     make(chan messages.StackResponse),
		compilerReqs:  make(chan compiler.CompilerResponse),
		logEventReqs:  make(chan map[string]string, 1),
	}

	home, _ := os.UserHomeDir()
	// Set shell and completers
	shellCompleter := session.getCompleter("main")
	session.Shell, _ = readline.NewEx(&readline.Config{
		HistoryFile:       home + "/.wiregost/client/.history",
		AutoComplete:      shellCompleter,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		// FilterInputRune: To be used later if needed
	})
	// Set Prompt
	session.prompt = NewPrompt(session)

	// Set Auth
	session.user = NewUser()
	session.user.LoadCreds()

	// Load saved servers
	session.LoadEndpointList()
	session.GetDefaultEndpoint()
	session.connected = false

	// Connect to default server
	session.Connect()

	// Launch console but give time to connect
	time.Sleep(time.Millisecond * 50)
	session.Start()

	return session
}

func (s *Session) Start() {

	// Eventually close the session
	defer s.Shell.Close()

	// Authenticate
	s.user.Authenticate()
	RefreshPrompt(s.prompt, s.Shell)

	// Read commands
	for {
		line, err := s.Shell.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cmd := strings.Fields(line)

		if len(cmd) > 0 {
			switch s.menuContext {
			case "main":
				s.mainMenuCommand(cmd)
			case "module":
				s.moduleMenuCommand(cmd)
			case "compiler":
				s.compilerMenuCommand(cmd)
			}
		}

		// Refresh shell & prompt after each command, at least.
		s.Refresh()
	}
}

func (s *Session) Refresh() {
	RefreshPrompt(s.prompt, s.Shell)
	s.Shell.Refresh()
}

func (s *Session) mainMenuCommand(cmd []string) {
	switch cmd[0] {
	// Core Commands
	case "help":
		helpHandler(cmd)
	case "cd":
		changeDirHandler(cmd)
	case "mode":
		mode := setModeHandler(cmd, s.Shell.IsVimMode())
		s.Shell.SetVimMode(mode)
	case "!":
		shellHandler(cmd[1:])
	case "exit":
		exit()
	case "set":
		s.SetOption(cmd)
	case "get":
		s.GetOption(cmd)
	// Endpoint
	case "endpoint":
		switch cmd[1] {
		case "list":
			s.ListEndpoints()
		case "add":
			s.AddEndpoint()
		case "connect":
			s.EndpointConnect(cmd)
		case "delete":
			s.DeleteEndpoint(cmd)
		}
	// Workspace
	case "workspace":
		switch cmd[1] {
		case "switch":
			s.WorkspaceSwitch(cmd)
		case "new":
			s.WorkspaceNew(cmd)
		case "delete":
			s.WorkspaceDelete(cmd)
		case "list":
			s.WorkspaceList(cmd)
		}
	case "log":
		switch cmd[1] {
		case "level":
			s.SetLogLevel(cmd)
		case "show":
			s.LogShow(cmd)
		}
	// Module
	case "use":
		s.UseModule(cmd)
	// Stack
	case "stack":
		switch len(cmd) {
		case 1:
			s.StackShow()
		case 2:
			switch cmd[1] {
			case "show":
				s.StackShow()
			case "pop":
				s.StackPop(cmd)
			}
		case 3:
			switch cmd[1] {
			case "use":
				s.StackUse(cmd)
			case "pop":
				s.StackPop(cmd)
			}
		}
	// Compiler
	case "compiler":
		s.UseCompiler()
	// Server
	case "server":
		switch cmd[1] {
		case "reload":
			s.ServerReload(cmd)
		case "start":
			s.ServerStart(cmd)
		case "stop":
			s.ServerStop(cmd)
		case "generate_certificate":
			s.GenerateCertificate(cmd)
		}
	}
}

func (s *Session) moduleMenuCommand(cmd []string) {
	switch cmd[0] {
	// Core Commands
	case "help":
		helpHandler(cmd)
	case "cd":
		changeDirHandler(cmd)
	case "mode":
		mode := setModeHandler(cmd, s.Shell.IsVimMode())
		s.Shell.SetVimMode(mode)
	case "!":
		shellHandler(cmd[1:])
	case "get":
		s.GetOption(cmd)
	case "exit":
		exit()
	// Endpoint
	case "endpoint":
		switch cmd[1] {
		case "list":
			s.ListEndpoints()
		case "add":
			s.AddEndpoint()
		case "connect":
			s.EndpointConnect(cmd)
		case "delete":
			s.DeleteEndpoint(cmd)
		}
	// Workspace
	case "workspace":
		switch cmd[1] {
		case "switch":
			s.WorkspaceSwitch(cmd)
		case "new":
			s.WorkspaceNew(cmd)
		case "list":
			s.WorkspaceList(cmd)
		}
	case "log":
		switch cmd[1] {
		case "level":
			s.SetLogLevel(cmd)
		case "show":
			s.LogShow(cmd)
		}
	// Module
	case "use":
		s.UseModule(cmd)
	case "show":
		switch cmd[1] {
		case "options":
			s.ShowOptions(cmd)
		case "info":
			s.ShowInfo()
		}
	case "info":
		s.ShowInfo()
	case "set":
		s.SetModuleOption(cmd)
	case "back":
		s.BackModule()
	// Stack
	case "stack":
		switch len(cmd) {
		case 1:
			s.StackShow()
		case 2:
			switch cmd[1] {
			case "show":
				s.StackShow()
			case "pop":
				s.StackPop(cmd)
			}
		case 3:
			switch cmd[1] {
			case "use":
				s.StackUse(cmd)
			case "pop":
				s.StackPop(cmd)
			}
		}
	// Server
	case "server":
		switch cmd[1] {
		case "reload":
			s.ServerReload(cmd)
		case "start":
			s.ServerStart(cmd)
		case "stop":
			s.ServerStop(cmd)
		case "generate_certificate":
			s.GenerateCertificate(cmd)
		}
	// Compiler
	case "compiler":
		s.UseCompiler()
		// Server
	}
}

func (s *Session) compilerMenuCommand(cmd []string) {
	switch cmd[0] {
	case "help":
		compilerHelp()
	case "back":
		s.QuitCompiler()
	case "list":
		switch cmd[1] {
		case "parameters":
			s.ShowCompilerOptions(cmd)
		}
	case "set":
		s.SetCompilerOption(cmd)
	}
}
