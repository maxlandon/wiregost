package cli

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/tui"
)

var (
	shellMenuContext = "main"
	user             *User
)

func Shell() {

	shellCompleter := getCompleter("main")
	promptString := NewPrompt()

	p, err := readline.NewEx(&readline.Config{
		Prompt:            "test",
		HistoryFile:       "/tmpt/testfile.tmp",
		AutoComplete:      shellCompleter,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		// FilterInputRune: To be used later if needed
	})

	if err != nil {
		fmt.Println(tui.Red("There was an error with the provided input"))
		fmt.Println(tui.Red(err.Error()))
	}

	// Authenticate
	user = NewUser()
	user.LoadCreds()
	user.Authenticate()

	// Set prompt
	prompt := p
	Refresh(promptString, prompt)
	defer prompt.Close()

	log.SetOutput(prompt.Stderr())

	// prompt.Config.AutoComplete = getCompleter("server")
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
			switch shellMenuContext {
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
				}

			}
		}

		// Refresh prompt after each command, at least
		Refresh(promptString, prompt)
	}
}

// Refresh prompt
func Refresh(prompt Prompt, input *readline.Instance) {
	p, _ := prompt.Render()
	_, m := prompt.Render()
	// p, _ := s.parseEnvTokens(s.Prompt.Render(s))
	fmt.Println()
	fmt.Println(p)
	input.SetPrompt(m)
	input.Refresh()
}
