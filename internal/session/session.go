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

// Session is the central object of a Wiregost client shell session.
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
	CurrentWorkspaceID int
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
	compilerReqs  chan compiler.Response
	logEventReqs  chan map[string]string
}

// NewSession instantiates a new Session object.
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
		compilerReqs:  make(chan compiler.Response),
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
	session.prompt = newPrompt(session)

	// Set Auth
	session.user = NewUser()
	session.user.LoadCreds()

	// Load saved servers
	session.loadEndpointList()
	session.getDefaultEndpoint()
	session.connected = false

	// Connect to default server
	session.connect()

	// Launch console but give time to connect
	time.Sleep(time.Millisecond * 50)
	session.start()

	return session
}

func (s *Session) start() {

	// Eventually close the session
	defer s.Shell.Close()

	// Authenticate
	s.user.Authenticate()
	refreshPrompt(s.prompt, s.Shell)

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
		s.refresh()
	}
}

func (s *Session) refresh() {
	refreshPrompt(s.prompt, s.Shell)
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
		s.setOption(cmd)
	case "get":
		s.setOption(cmd)
	// Endpoint
	case "endpoint":
		switch cmd[1] {
		case "list":
			s.listEndpoints()
		case "add":
			s.addEndpoint()
		case "connect":
			s.endpointConnect(cmd)
		case "delete":
			s.deleteEndpoint(cmd)
		}
	// Workspace
	case "workspace":
		switch cmd[1] {
		case "switch":
			s.workspaceSwitch(cmd)
		case "new":
			s.workspaceNew(cmd)
		case "delete":
			s.workspaceDelete(cmd)
		case "list":
			s.workspaceList(cmd)
		}
	case "log":
		switch cmd[1] {
		case "level":
			s.setLogLevel(cmd)
		case "show":
			s.logShow(cmd)
		}
	// Module
	case "use":
		s.useModule(cmd)
	// Stack
	case "stack":
		switch len(cmd) {
		case 1:
			s.stackShow()
		case 2:
			switch cmd[1] {
			case "show":
				s.stackShow()
			case "pop":
				s.stackPop(cmd)
			}
		case 3:
			switch cmd[1] {
			case "use":
				s.stackUse(cmd)
			case "pop":
				s.stackPop(cmd)
			}
		}
	// Compiler
	case "compiler":
		s.useCompiler()
	// Server
	case "server":
		switch cmd[1] {
		case "reload":
			s.serverReload(cmd)
		case "start":
			s.serverStart(cmd)
		case "stop":
			s.serverStop(cmd)
		case "generate_certificate":
			s.generateCertificate(cmd)
		case "list":
			s.serverList(cmd)
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
		s.getOption(cmd)
	case "exit":
		exit()
	// Endpoint
	case "endpoint":
		switch cmd[1] {
		case "list":
			s.listEndpoints()
		case "add":
			s.addEndpoint()
		case "connect":
			s.endpointConnect(cmd)
		case "delete":
			s.deleteEndpoint(cmd)
		}
	// Workspace
	case "workspace":
		switch cmd[1] {
		case "switch":
			s.workspaceSwitch(cmd)
		case "new":
			s.workspaceNew(cmd)
		case "list":
			s.workspaceList(cmd)
		}
	case "log":
		switch cmd[1] {
		case "level":
			s.setLogLevel(cmd)
		case "show":
			s.logShow(cmd)
		}
	// Module
	case "use":
		s.useModule(cmd)
	case "show":
		switch cmd[1] {
		case "options":
			s.showOptions(cmd)
		case "info":
			s.showInfo()
		}
	case "info":
		s.showInfo()
	case "set":
		s.setModuleOption(cmd)
	case "back":
		s.backModule()
	// Stack
	case "stack":
		switch len(cmd) {
		case 1:
			s.stackShow()
		case 2:
			switch cmd[1] {
			case "show":
				s.stackShow()
			case "pop":
				s.stackPop(cmd)
			}
		case 3:
			switch cmd[1] {
			case "use":
				s.stackUse(cmd)
			case "pop":
				s.stackPop(cmd)
			}
		}
	// Server
	case "server":
		switch cmd[1] {
		case "reload":
			s.serverReload(cmd)
		case "start":
			s.serverStart(cmd)
		case "stop":
			s.serverStop(cmd)
		case "generate_certificate":
			s.generateCertificate(cmd)
		case "list":
			s.serverList(cmd)
		}
	// Compiler
	case "compiler":
		s.useCompiler()
		// Server
	}
}

func (s *Session) compilerMenuCommand(cmd []string) {
	switch cmd[0] {
	case "help":
		compilerHelp()
	case "back":
		s.quitCompiler()
	case "list":
		switch cmd[1] {
		case "parameters":
			s.showCompilerOptions(cmd)
		}
	case "set":
		s.setCompilerOption(cmd)
	}
}
