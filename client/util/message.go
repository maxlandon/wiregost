package util

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

var (
	// Info - All normal message
	Info = fmt.Sprintf("%s[-]%s ", tui.BLUE, tui.RESET)
	// Warn - Errors in parameters, notifiable events in modules/sessions
	Warn = fmt.Sprintf("%s[!]%s ", tui.YELLOW, tui.RESET)
	// Error - Error in commands, filters, modules and implants.
	Error = fmt.Sprintf("%s[!]%s ", tui.RED, tui.RESET)
	// Success - Success events
	Success = fmt.Sprintf("%s[*]%s ", tui.GREEN, tui.RESET)

	// Infof - formatted
	Infof = fmt.Sprintf("%s[-] ", tui.BLUE)
	// Warnf - formatted
	Warnf = fmt.Sprintf("%s[!] ", tui.YELLOW)
	// Errorf - formatted
	Errorf = fmt.Sprintf("%s[!] ", tui.RED)
	// Sucessf - formatted
	Sucessf = fmt.Sprintf("%s[*] ", tui.GREEN)

	//RPCError - Errors from the server
	RPCError = fmt.Sprintf("%s[RPC Error]%s ", tui.RED, tui.RESET)
	// CommandError - Command input error
	CommandError = fmt.Sprintf("%s[Command Error]%s ", tui.RED, tui.RESET)
	// ParserError - Failed to parse some tokens in the input
	ParserError = fmt.Sprintf("%s[Parser Error]%s ", tui.RED, tui.RESET)
	// DBError - Data Service error
	DBError = fmt.Sprintf("%s[DB Error]%s ", tui.RED, tui.RESET)
)
