package cli

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

// MODULE HANDLERS
//---------------------------------------------------------------------------

// Function used for description paragraphs
func wrap(text string, lineWidth int) (wrapped string) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}

func (s *Session) UseModule(cmd []string) {
	s.Send(cmd)
	mod := <-moduleReqs
	// Switch shell context
	s.shell.Config.AutoComplete = s.getCompleter("module")
	s.shellMenuContext = "module"
	s.moduleContext = mod.ModuleName
	// Switch prompt context
	CurrentModule = mod.ModuleName
	// Add code to change current module in the prompt
}

func (s *Session) ShowOptions() {
	s.Send(strings.Fields("show options"))
	mod := <-moduleReqs
	m := mod.Modules[0]

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(3, 50)
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
	)
	table.SetBorder(false)
	// TODO add option for agent alias here
	table.Append([]string{"Agent", m.Agent.String(), "true", "Agent on which to run module " + m.Name})
	for _, v := range m.Options {
		table.Append([]string{v.Name, v.Value, strconv.FormatBool(v.Required), v.Description})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) ShowInfo() {
	s.Send(strings.Fields("show options"))
	mod := <-moduleReqs
	m := mod.Modules[0]

	// Info
	fmt.Printf("%sModule:%s\r\n\t%s\r\n", tui.YELLOW, tui.RESET, m.Name)
	fmt.Printf("%sPlatform:%s\r\n\t%s\\%s\\%s\r\n", tui.YELLOW, tui.RESET, m.Platform, m.Arch, m.Lang)
	fmt.Printf("%sModule Authors:%s\n", tui.YELLOW, tui.RESET)
	for a := range m.Author {
		fmt.Printf("\t%s\n", m.Author[a])
	}
	fmt.Printf("%sCredits:%s\n", tui.YELLOW, tui.RESET)
	for c := range m.Credits {
		fmt.Printf("\t%s\n", m.Credits[c])
	}
	fmt.Printf("%sDescription:%s\r\n", tui.YELLOW, tui.RESET)
	fmt.Println(wrap(m.Description, 140))
	fmt.Println()
	// Table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(3, 50)
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
	)
	table.SetBorder(false)
	// TODO add option for agent alias here
	table.Append([]string{"Agent", m.Agent.String(), "true", "Agent on which to run module " + m.Name})
	for _, v := range m.Options {
		table.Append([]string{v.Name, v.Value, strconv.FormatBool(v.Required), v.Description})
	}
	fmt.Println()
	table.Render()
	fmt.Println()
	fmt.Printf("%sNotes:%s\n", tui.YELLOW, tui.RESET)
	fmt.Println(wrap(m.Notes, 140))
	fmt.Println()
}

func (s *Session) GetModuleList(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs

	list := mod.Modules
	fmt.Println(list)
}

func (s *Session) SetModuleOption(cmd []string) {
	s.Send(cmd)
	opt := <-moduleReqs
	if opt.Status != "" {
		fmt.Println()
		fmt.Println(opt.Status)
	}
	if opt.Error != "" {
		fmt.Println()
		fmt.Println(opt.Error)
	}
}

func (s *Session) SetAgent(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs
	fmt.Println(mod)
	// Add some verification that agent is correctly set here.
}

func (s *Session) RunModule(cmd []string) {
	// Send(cmd)
	mod := <-moduleReqs
	fmt.Println(mod)
	// Add some verification that agent is correctly set here.
}

func (s *Session) BackModule() {
	s.shell.Config.AutoComplete = s.getCompleter("main")
	s.shellMenuContext = "main"
	s.moduleContext = ""
	CurrentModule = ""
}

// AGENT HANDLERS
//---------------------------------------------------------------------------
func (s *Session) ListAgents(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion for listing agents
}

func (s *Session) InfoAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Add func (s *Session)tion to print info.
}

func (s *Session) RemoveAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) DownloadAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) UploadAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) ExecuteShellCodeAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) KillAgent(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) SetAgentOption(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Check to see if we need an answer here, or if we can just go on
	// and wait for it via the log.
}

func (s *Session) GetAgentShell(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Don't know how we will handle this one
}

func (s *Session) BackMainMenu(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Handle change of state here
}

func (s *Session) MainMenu(cmd []string) {
	// Send(cmd)
	agent := <-agentReqs
	fmt.Println(agent)
	// Handle change of state here
}

// LOG HANDLERS
//---------------------------------------------------------------------------

func (s *Session) LogLevel(cmd []string) {
	// Send(cmd)
	log := <-logReqs
	fmt.Println(log)
	// Handle change of state here
}

func (s *Session) LogShow(cmd []string) {
	// Send(cmd)
	log := <-logReqs
	fmt.Println(log)
	// Handle printing the logs here
}

// Handle all log messages coming from the server
func (s *Session) LogListen() {
	go func() {
		for {
			msg := <-logReqs
			fmt.Println(msg)
		}
	}()
}

// WORKSPACE HANDLERS
//---------------------------------------------------------------------------

func (s *Session) WorkspaceList(cmd []string) {
	s.Send(cmd)
	workspace := <-workspaceReqs
	fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceSwitch(cmd []string) {
	s.currentWorkspace = cmd[2]
	CurrentWorkspace = cmd[2]
	// fmt.Println(CurrentWorkspace)
	s.Send(cmd)
	workspace := <-workspaceReqs
	s.CurrentWorkspaceId = workspace.WorkspaceId
	// fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceDelete(cmd []string) {
	// Send(cmd)
	workspace := <-workspaceReqs
	fmt.Println(workspace)
	// Handle change of state here
}

func (s *Session) WorkspaceNew(cmd []string) {
	s.Send(cmd)
	// workspace := <-workspaceReqs
	// fmt.Println(workspace)
	// Handle change of state here
}

// STACK HANDLERS
//---------------------------------------------------------------------------

func (s *Session) StackShow() {
	s.Send(strings.Fields("stack show"))
	stack := <-moduleReqs

	// Print stack
	fmt.Println(tui.Dim("The stack stores a list of previously loaded modules and their state (options, agents) "))
	fmt.Println(tui.Dim("Source local scripts are in /data/src/."))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(1, 50)
	table.SetHeader([]string{"Name", "Source Local", "Language"})
	table.SetAutoWrapText(true)
	table.SetReflowDuringAutoWrap(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
	)
	table.SetBorder(false)
	// TODO add option for agent alias here
	for i := len(stack.Modules) - 1; i >= 0; i-- {
		if strings.ToLower(strings.TrimSuffix(strings.Join(stack.Modules[i].Path, "/"), ".json")) == strings.ToLower(s.moduleContext) {
			table.Rich([]string{stack.Modules[i].Name, strings.TrimPrefix(strings.Join(stack.Modules[i].SourceLocal, "/"), "data/src"), stack.Modules[i].Lang},
				[]tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				})
		} else {
			table.Append([]string{stack.Modules[i].Name, strings.TrimPrefix(strings.Join(stack.Modules[i].SourceLocal, "/"), "data/src"), stack.Modules[i].Lang})
		}
	}
	fmt.Println()
	table.Render()
}

func (s *Session) StackPop(cmd []string) {
	s.Send(cmd)
	switch len(cmd) {
	case 2:
		s.moduleContext = ""
		CurrentModule = ""
	case 3:
		if strings.ToLower(cmd[2]) == strings.ToLower(s.moduleContext) {
			s.moduleContext = ""
			CurrentModule = ""
		}
	}
	// Temporary: return to main menu completion.
	// This will change when the code will handle fallback on next module in stack.
	s.shell.Config.AutoComplete = s.getCompleter("main")

}

// ENDPOINT HANDLERS
//---------------------------------------------------------------------------

func (s *Session) EndpointConnect() { // This command will be used to send connection demands to a server.
	// Send(cmd)
	endpoint := <-endpointReqs
	fmt.Println(endpoint)
}

// SERVER HANDLERS
//---------------------------------------------------------------------------
func (s *Session) ServerReload() {
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
		Context:            s.shellMenuContext,
		CurrentModule:      s.moduleContext,
		Command:            []string{"server", "reload"},
		ServerParams:       params,
	}
	enc := json.NewEncoder(writer)
	err := enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	status := <-serverReqs
	fmt.Println()
	fmt.Println(status.Status)
}
