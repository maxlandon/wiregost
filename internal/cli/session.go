package cli

import (
	"io"
	"log"
	"strings"

	// 3rd party
	"github.com/chzyer/readline"
)

var ()

type Session struct {
	shell              *readline.Instance
	prompt             Prompt
	user               *User
	shellMenuContext   string
	moduleContext      string
	workspaceContext   string
	currentWorkspace   string
	CurrentWorkspaceId int
}

func NewSession() *Session {
	session := &Session{}

	shellCompleter := session.getCompleter("main")

	session.shell, _ = readline.NewEx(&readline.Config{
		HistoryFile:       "/tmpt/testfile.tmp",
		AutoComplete:      shellCompleter,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		// FilterInputRune: To be used later if needed
	})

	// Set Prompt
	session.prompt = NewPrompt()

	// Set Auth
	session.user = NewUser()
	session.user.LoadCreds()

	// Set Context
	session.shellMenuContext = "main"

	// Connect to default server
	Connect()

	// Launch console
	session.Shell()

	return session
}

func (session *Session) Shell() {

	// Authenticate
	session.user.Authenticate()

	// Set prompt
	prompt := session.shell
	Refresh(session.prompt, prompt)
	defer prompt.Close()

	log.SetOutput(prompt.Stderr())

	// Read commands
	for {
		line, err := prompt.Readline()
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
			switch session.shellMenuContext {
			case "main":
				switch cmd[0] {
				// Core Commands
				case "help":
					helpHandler(cmd)
				case "cd":
					changeDirHandler(cmd)
				case "mode":
					mode := setModeHandler(cmd, prompt.IsVimMode())
					prompt.SetVimMode(mode)
				case "!":
					shellHandler(cmd[1:])
				case "exit":
					exit()
				// Workspace
				case "workspace":
					switch cmd[1] {
					case "switch":
						session.WorkspaceSwitch(cmd)
					case "new":
						session.WorkspaceNew(cmd)
					}
				// Module
				case "use":
					session.UseModule(cmd)
				case "show":
					switch cmd[1] {
					case "options":
						session.ShowOptions()
					case "info":
						session.ShowInfo()
					}
				case "info":
					session.ShowInfo()
				case "set":
					session.SetModuleOption(cmd)
				// Stack
				case "stack":
					switch len(cmd) {
					case 1:
						session.StackShow()
					case 2:
						switch cmd[1] {
						case "show":
							session.StackShow()
						case "pop":
							session.StackPop(cmd)
						}
					case 3:
						session.StackPop(cmd)
					}

				}
			}
		}

		// Refresh prompt after each command, at least
		Refresh(session.prompt, prompt)
	}
}
