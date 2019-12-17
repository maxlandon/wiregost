package session

// The session package contains the code for the WireGost client shell, ghost.
// The structure is similar to the one of Bettercap, the wireless toolkit written
// in Go.

// This file contains the Session object, which is central to the shell.
// It interfaces with all other components needed by the shell, such as prompts,
// command handlers, client configuration, etc.

// The most basic functions needed by the shell are also here, such as New(),
// Refresh(), ReadLine() and Run().

import (
	"fmt"
	"os"
	"time"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

var userHomeDir, err = os.UserHomeDir()
var HistoryFile = userHomeDir + "/.wiregost/client/.history"

type UnknownCommandCallback func(cmd string) bool

type Session struct {
	StartedAt       time.Time
	Active          bool
	Prompt          Prompt
	Input           *readline.Instance
	CommandHandlers []CommandHandler
	Config          *Config
	User            *User
	ServerManager   *ServerManager
	CurrentDir      string
	UnkCmdCallback  UnknownCommandCallback
}

// Instantiates a new Session object
func New() (*Session, error) {

	s := &Session{
		Prompt: NewPrompt(),
		Config: NewConfig(),
		User:   NewUser(),

		CommandHandlers: make([]CommandHandler, 0),
		UnkCmdCallback:  nil,
	}
	// Load User credentials
	s.User.LoadCreds()

	// Start Server Manager
	s.ServerManager = NewServerManager(s.User)

	// Register all command handlers
	s.registerCoreHandlers()
	s.registerConfigHandlers()
	s.registerHelpHandlers()
	s.registerHistoryHandlers()
	s.registerServerHandlers()

	return s, nil
}

// Starts a Session based on a Session object
func (s *Session) Start() (err error) {

	if err := s.setupReadline(); err != nil {
		return err
	}

	// Loading all configuration elements MIGHT BE IMPORTANT IN THE FUTURE, FOR OTHER SERVICES/FUNCTIONS !!!!!!!!
	// s.Config.LoadConfig()

	// Load User Creds and authenticate
	s.User.Authenticate()

	s.StartedAt = time.Now()
	s.Active = true

	return err
}

// Not sure this will be useful, taken from Bettercap
func (s *Session) Lock() {

}

// Not sure this will be useful, taken from Bettercap
func (s *Session) Unlock() {

}

// Close the Session
func (s *Session) Close() {

}

// Refreshes the console each time it is needed.
func (s *Session) Refresh() {
	p, _ := s.Prompt.Render(s)
	_, m := s.Prompt.Render(s)
	// p, _ := s.parseEnvTokens(s.Prompt.Render(s))
	fmt.Println()
	fmt.Println(p)
	s.Input.SetPrompt(m)
	s.Input.Refresh()
}

func (s *Session) ReadLine() (string, error) {
	s.Refresh()
	return s.Input.Readline()
}

// Not sure this will be used here. Original parses and make primitive command dispatch.
func (s *Session) Run(line string) error {
	line = str.TrimRight(line)
	// line = reCmdSpaceCleaner.ReplaceAllString(line, "$1 $2")

	// is it a core command?
	for _, h := range s.CommandHandlers {
		if parsed, args := h.Parse(line); parsed {
			return h.Exec(args, s)
		}
	}

	// is it a proxy module custom command?
	if s.UnkCmdCallback != nil && s.UnkCmdCallback(line) {
		return nil
	}

	// If command is not valid
	return fmt.Errorf("unknown or invalid syntax \"%s%s%s\", type %shelp%s for the help menu.",
		tui.BOLD, line, tui.RESET, tui.BOLD, tui.RESET)
}
