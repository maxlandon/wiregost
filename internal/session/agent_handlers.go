package session

import (
	"fmt"
	"os"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) listAgents(cmd []string) {
	s.send([]string{"agent", "show"})
	agents := <-s.agentReqs

	// Render list
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetHeader([]string{"Agent ID", "Platform", "UserName", "HostName", "Proto", "Status CheckIn"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	table.SetBorder(false)

	for _, a := range agents.Infos {
		// Convert proto (i.e. h2 or hq) to user friendly string
		var proto string
		if a["protocol"] == "https" {
			proto = "HTTP/1.1 (https)"
		} else if a["protocol"] == "h2" {
			proto = "HTTP/2 (h2)"
		} else if a["protocol"] == "hq" {
			proto = "QUIC (hq)"
		}
		table.Append([]string{a["id"], a["platform"] + "/" + a["arch"], a["username"], a["hostname"], proto, a["statusCheckIn"]})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) infoAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion to print info.
}

func (s *Session) removeAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) downloadAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) uploadAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) executeShellCodeAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) killAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) setAgentOption(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) getAgentShell(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) backMainMenu(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Handle change of state here
}

func (s *Session) mainMenu(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Handle change of state here
}
