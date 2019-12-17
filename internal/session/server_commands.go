package session

import (
	"os"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

// This file contains all server command handlers and their registering function.

// Connect User
func (s *Session) ConnectUserHandler(args []string, sess *Session) error {
	filter := ""
	if len(args) == 2 {
		filter = str.Trim(args[1])
	}

	for _, l := range s.ServerManager.SavedServers {
		str := l.FQDN + tui.Dim(" at ") + l.IPAddress + ":" + strconv.Itoa(l.Port)
		if filter == str {
			s.ServerManager.ConnectToServer(s.User, l)
			return nil
		}
	}

	s.ServerManager.ConnectToServer(s.User, s.ServerManager.CurrentServer)
	return nil
}

// List Servers
func (s *Session) ListServersHandler(args []string, sess *Session) error {
	columns := []string{
		tui.Yellow("FQDN (Common Name)"),
		tui.Yellow("Address"),
		tui.Yellow("Certificate"),
		tui.Yellow("Connected"),
		tui.Yellow("Default"),
	}

	rows := [][]string{}

	for _, l := range s.ServerManager.SavedServers {
		row := []string{}
		// Name
		row = append(row, l.FQDN)
		// IP:Port
		address := l.IPAddress + ":" + strconv.Itoa(l.Port)
		row = append(row, address)
		// Certificate name (removing path)
		row = append(row, l.Certificate)
		// Connected
		if s.ServerManager.CurrentServer == l {
			row = append(row, tui.Green("Connected"))
		}
		if s.ServerManager.CurrentServer != l {
			row = append(row, " ")
		}
		// Default
		if l.IsDefault == true {
			row = append(row, "default")
		}
		if l.IsDefault == false {
			row = append(row, " ")
		}
		//Append to servers list
		rows = append(rows, row)
	}
	// Print table
	tui.Table(os.Stdout, columns, rows)
	return nil
}

// Register all handlers defined above
func (s *Session) registerServerHandlers() {
	//Register User
	s.addHandler(NewCommandHandler("server.connect",
		"^(server.connect|\\?)(.*)$",
		"Connect to the specified server",
		s.ConnectUserHandler),
		readline.PcItem("server.connect", readline.PcItemDynamic(func(prefix string) []string {
			prefix = str.Trim(prefix[14:])
			servers := []string{}
			for _, l := range s.ServerManager.SavedServers {
				server := l.FQDN + tui.Dim(" at ") + l.IPAddress + ":" + strconv.Itoa(l.Port)
				servers = append(servers, server)
			}
			return servers
		})))

	// List servers
	s.addHandler(NewCommandHandler("server.list",
		"server.list",
		"List available servers and their settings",
		s.ListServersHandler),
		readline.PcItem("server.list"))
}
