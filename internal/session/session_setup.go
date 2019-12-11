package session

// This file contains all functions that the shell needs for its own
// setup.

// This includes registering all the completers and setting things
// such as command history saving behavior.

import (
	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/fs"
)

func (s *Session) setupReadline() (err error) {
	prefixCompleters := make([]readline.PrefixCompleterInterface, 0)
	for _, h := range s.CommandHandlers {
		if h.Completer == nil {
			prefixCompleters = append(prefixCompleters, readline.PcItem(h.Name))
		} else {
			prefixCompleters = append(prefixCompleters, h.Completer)
		}
	}

	tree := make(map[string][]string)
	// In Bettercap, there is a bit of code here that we haven't find useful
	// to use yet.

	for root, subElems := range tree {
		item := readline.PcItem(root)
		item.Children = []readline.PrefixCompleterInterface{}
		for _, child := range subElems {
			item.Children = append(item.Children, readline.PcItem(child))
		}
		prefixCompleters = append(prefixCompleters, item)
	}

	history, _ := fs.Expand(HistoryFile)

	cfg := readline.Config{
		HistoryFile:            history,
		DisableAutoSaveHistory: false,
		HistoryLimit:           2000,
		InterruptPrompt:        "^C",
		EOFPrompt:              "^D",
		AutoComplete:           readline.NewPrefixCompleter(prefixCompleters...),
	}

	s.Input, err = readline.NewEx(&cfg)
	return err
}
