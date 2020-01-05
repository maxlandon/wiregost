package session

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) ServerStart(cmd []string) {
	s.Send(cmd)
	status := <-s.serverReqs
	fmt.Println()
	fmt.Println(status.Status)
}

func (s *Session) ServerReload(cmd []string) {
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
		CurrentWorkspaceId: s.CurrentWorkspaceId,
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

func (s *Session) ServerStop(cmd []string) {
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
		CurrentWorkspaceId: s.CurrentWorkspaceId,
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

func (s *Session) ServerList(cmd []string) {
	s.Send(cmd)
	serv := <-s.serverReqs

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(5, 50)
	table.SetHeader([]string{"Workspace", "Address", "Protocol", "State", "PSK", "Certificate"})
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

	running := ""
	for _, v := range serv.ServerList {
		if v["state"] == "true" {
			running = tui.Green("Running")
		} else {
			running = tui.Red("Stopped")
		}
		table.Append([]string{v["workspace"], v["address"], v["protocol"], running, v["psk"], v["certificate"]})
	}
	fmt.Println()
	table.Render()

}

func (s *Session) GenerateCertificate(cmd []string) {
	s.Send(cmd)
	server := <-s.serverReqs
	fmt.Println()
	fmt.Println(server.Status)
}
