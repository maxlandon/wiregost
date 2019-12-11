package session

// This file contains all history command handlers, and their registering function.

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

// Show last x nb of commands in history
func (s *Session) historyShowHandler(nb int) error {

	// Determine history length for subsequent selection
	hist, _ := os.Open(s.Config.HistoryFile)
	hlength := 0
	scanner := bufio.NewScanner(hist)
	for scanner.Scan() {
		hlength += 1
	}
	hist.Close()

	// Read history and output last x commands
	file, _ := os.Open(s.Config.HistoryFile)
	defer file.Close()
	count := 1
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		if hlength <= nb {
			line := tui.Dim(strconv.Itoa(count)) + tui.Dim(") ") + scan.Text()
			fmt.Println(line)
			count += 1
		}
		hlength -= 1
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("%sError parsing history file. %s\n", tui.RED, tui.RESET)
		log.Fatal(err)
	}

	return nil
}

// General history handler
func (s *Session) historyHandler(args []string, sess *Session) error {
	// Now, only one command is available for history: show it.
	filter := ""
	nb := 0
	if len(args) == 2 {
		filter = str.Trim(args[1])
		nb, _ = strconv.Atoi(filter)
	}

	s.historyShowHandler(nb)

	return nil
}

// Register all history handlers
func (s *Session) registerHistoryHandlers() {
	// Exit
	s.addHandler(NewCommandHandler("history.show",
		"^(history.show|\\?)(.*)$",
		"Display last x number of history commands.",
		s.historyHandler),
		readline.PcItem("history.show"))
}
