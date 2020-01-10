package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/internal/messages"
	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
)

func (s *Session) workspaceList(cmd []string) {
	s.send(cmd)
	workspace := <-s.workspaceReqs

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(1, 50)
	table.SetHeader([]string{"Name", "Description", "Boundary"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	table.SetBorder(false)
	for _, w := range workspace.WorkspaceInfos {
		table.Append([]string{w[0], w[1], w[2]})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) workspaceSwitch(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: provide a workspace name'\n", tui.RED, tui.RESET)
		return
	}

	// Check workspace, so we don't send garbage to the server.
	s.send([]string{"workspace", "list"})
	workspaceList := <-s.workspaceReqs
	recognized := false
	for _, w := range workspaceList.WorkspaceInfos {
		if w[0] == cmd[2] {
			recognized = true
		}
	}
	if !recognized {
		fmt.Printf("%s[!]%s Error: workspace %s does not exist.'\n", tui.RED, tui.RESET, cmd[2])
		return
	}
	// Eventually send command
	s.send(cmd)
	workspace := <-s.workspaceReqs
	server := <-s.serverReqs
	fmt.Println()
	fmt.Printf(workspace.Result)
	fmt.Println(server.Status)
	// Change menu context and completion
	s.menuContext = "main"
	s.Shell.Config.AutoComplete = s.getCompleter("main")
	// Change workspace variables
	s.CurrentWorkspaceID = workspace.WorkspaceID
	s.currentWorkspace = cmd[2]
	// Reset modules
	s.currentModule = ""
	// Change server variables
	s.currentServerID = server.ServerID
	s.serverRunning = server.ServerRunning
	// Reset agent variables
	s.currentAgentID = uuid.FromStringOrNil("")
}

func (s *Session) workspaceDelete(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: provide a workspace name'\n", tui.RED, tui.RESET)
		return
	}

	// Check workspace, so we don't send garbage to the server.
	s.send([]string{"workspace", "list"})
	workspaceList := <-s.workspaceReqs
	recognized := false
	for _, w := range workspaceList.WorkspaceInfos {
		if w[0] == cmd[2] {
			recognized = true
		}
	}
	if !recognized {
		fmt.Printf("%s[!]%s Error: workspace %s does not exist.'\n", tui.RED, tui.RESET, cmd[2])
		return
	}

	// Check if workspace is not current one
	if cmd[2] == s.currentWorkspace {
		fmt.Println()
		fmt.Printf("%s[!]%s Cannot delete current workspace", tui.RED, tui.RESET)
		fmt.Println()
	} else {
		fmt.Println("Warning: deleting workspace will erase all agents and data related to it. Continue ? (Yes/No)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if (text == "yes\n") || (text == "Yes\n") {
			s.send(cmd)
			workspace := <-s.workspaceReqs
			fmt.Println()
			fmt.Println(workspace.Result)
		} else {
			fmt.Println()
			fmt.Println("[*] Aborted workspace deletion")
		}
	}
}

func (s *Session) workspaceNew(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: provide a workspace name'\n", tui.RED, tui.RESET)
		return
	}
	// send params if they are set
	params := make(map[string]string)
	for k, v := range s.Env {
		if strings.HasPrefix(k, "workspace") {
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
		WorkspaceParams:    params,
	}
	enc := json.NewEncoder(s.writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	s.writer.Flush()
	workspace := <-s.workspaceReqs
	fmt.Println()
	fmt.Println(workspace.Result)
}
