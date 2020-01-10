package session

import (
	// Standard
	"fmt"
	"os"

	// 3rd party
	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
)

func (s *Session) agentInteract(cmd []string) {
	// Check command
	if cmd[0] == "interact" {
		if len(cmd) < 2 {
			fmt.Printf("%s[!]%s Invalid command: select agent.\n", tui.RED, tui.RESET)
			return

		}
		s.menuContext = "agent"
		s.currentAgentID, _ = uuid.FromString(cmd[1])

	}
	if cmd[0] == "agent" {
		if len(cmd) < 3 {
			fmt.Printf("%s[!]%s Invalid command: select agent.\n", tui.RED, tui.RESET)
			return
		}
		s.menuContext = "agent"
		s.currentAgentID, _ = uuid.FromString(cmd[2])

	}
	// Change menu
	s.Shell.Config.AutoComplete = s.getCompleter("agent")
}

func (s *Session) listAgents(cmd []string) {
	s.send([]string{"agent", "show"})
	agents := <-s.agentReqs

	// Render list
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetHeader([]string{"Agent ID", "Platform", "UserName", "HostName", "Status", "Transport", "Status CheckIn"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
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
		table.Append([]string{a["id"], a["platform"] + "/" + a["arch"], a["username"], a["hostname"], a["status"], proto, a["statusCheckIn"]})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) infoAgent(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(agent.AgentInfo)

	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" Agent Information ")))
	fmt.Println()
	table.Render()
}

func (s *Session) listAgentDirectories(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) changeAgentDirectory(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) printAgentDirectory(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) agentCmd(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) downloadAgent(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) uploadAgent(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) executeShellCodeAgent(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) killAgent(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
	// Reset current agent
	s.currentAgentID = uuid.FromStringOrNil("")
	// Switch shell context
	if s.currentModule != "" {
		s.Shell.Config.AutoComplete = s.getCompleter("module")
		s.menuContext = "module"
		// Switch prompt context
	} else {
		s.Shell.Config.AutoComplete = s.getCompleter("main")
		s.menuContext = "main"
	}
}

func (s *Session) setAgentOption(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) getAgentShell(cmd []string) {
	s.send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent.Status)
}

func (s *Session) removeAgent(cmd []string) {
	// Send(cmd)
	agent := <-s.agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) backAgentMenu(cmd []string) {
	// Switch shell context
	if s.currentModule != "" {
		s.Shell.Config.AutoComplete = s.getCompleter("module")
		s.menuContext = "module"
		// Switch prompt context
	} else {
		s.Shell.Config.AutoComplete = s.getCompleter("main")
		s.menuContext = "main"
	}
}
