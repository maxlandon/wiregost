package session

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) serverStart(cmd []string) {
	s.send(cmd)
	status := <-s.serverReqs
	fmt.Println()
	fmt.Println(status.Status)
}

func (s *Session) serverReload(cmd []string) {
	// Fill up required parameters
	params := make(map[string]string)
	for k, v := range s.Env {
		if strings.HasPrefix(k, "server") {
			params[k] = v
		}
	}

	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceID: s.CurrentWorkspaceID,
		Context:            s.menuContext,
		CurrentModule:      s.currentModule,
		Command:            cmd,
		ServerParams:       params,
	}
	enc := json.NewEncoder(s.writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	s.writer.Flush()

	status := <-s.serverReqs
	fmt.Println()
	fmt.Println(status.Status)
}

func (s *Session) serverStop(cmd []string) {
	// Fill up required parameters
	params := make(map[string]string)
	for k, v := range s.Env {
		if strings.HasPrefix(k, "server") {
			params[k] = v
		}
	}

	// Used to fill the reload() function called after deleting
	// the server, so that it is ready to run again, with new
	// parameters
	msg := messages.ClientRequest{
		UserName:           s.user.Name,
		UserPassword:       s.user.PasswordHashString,
		CurrentWorkspace:   s.currentWorkspace,
		CurrentWorkspaceID: s.CurrentWorkspaceID,
		Context:            s.menuContext,
		CurrentModule:      s.currentModule,
		Command:            cmd,
		ServerParams:       params,
	}
	enc := json.NewEncoder(s.writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	s.writer.Flush()

	status := <-s.serverReqs
	fmt.Println()
	fmt.Println(status.Status)
}

func (s *Session) serverList(cmd []string) {
	// Get Servers
	s.send(cmd)
	serv := <-s.serverReqs

	// Get number of agents per server
	agents := make(map[string]int)
	for _, v := range serv.ServerList {
		s.send([]string{"agent", "list", v["id"]})
		a := <-s.agentReqs
		agents[v["id"]] = a.AgentNb[v["id"]]
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(6, 50)
	table.SetHeader([]string{"Workspace", "Address", "Protocol", "State", "PSK", "Agents", "Certificate"})
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

	running := ""
	for _, v := range serv.ServerList {
		if v["state"] == "true" {
			running = tui.Green("Running")
		} else {
			running = tui.Red("Stopped")
		}
		table.Append([]string{v["workspace"], v["address"], v["protocol"], running, v["psk"], strconv.Itoa(agents[v["id"]]), v["certificate"]})
	}
	fmt.Println()
	table.Render()

}

func (s *Session) generateCertificate(cmd []string) {
	s.send(cmd)
	fmt.Println()
	fmt.Println("  Generating Certificate and private key...")
	server := <-s.serverReqs
	fmt.Println(server.Status)
}
