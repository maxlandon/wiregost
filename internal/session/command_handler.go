package session

// All specific command handlers need to use a generic handler interface for sending commands.
// This file contains all the code related to this generic handler (parsing, adding and executing).

// CommandHandler is a generic handler that will be implemented by all commands

import (
	"regexp"
	"sync"

	"github.com/chzyer/readline"
)

type CommandHandler struct {
	sync.Mutex
	Name        string
	Description string
	Completer   *readline.PrefixCompleter
	Parser      *regexp.Regexp
	exec        func(args []string, s *Session) error
}

func NewCommandHandler(name string, expr string, desc string, exec func(args []string, s *Session) error) CommandHandler {
	return CommandHandler{
		Name:        name,
		Description: desc,
		Completer:   nil,
		Parser:      regexp.MustCompile(expr),
		exec:        exec,
	}
}

func (h *CommandHandler) Parse(line string) (bool, []string) {
	result := h.Parser.FindStringSubmatch(line)
	if len(result) == h.Parser.NumSubexp()+1 {
		return true, result[1:]
	} else {
		return false, nil
	}
}

func (h *CommandHandler) Exec(args []string, s *Session) error {
	h.Lock()
	defer h.Unlock()
	return h.exec(args, s)
}

// Add handler to session object. This is used each time a handler is registered
func (s *Session) addHandler(h CommandHandler, c *readline.PrefixCompleter) {
	h.Completer = c
	s.CommandHandlers = append(s.CommandHandlers, h)
}
