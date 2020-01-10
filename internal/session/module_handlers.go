package session

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) useModule(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: select module to use. \n", tui.RED, tui.RESET)
		return
	}
	// Check workspace, so we don't send garbage to the server.
	s.send([]string{"module", "list"})
	resp := <-s.moduleReqs
	list := resp.ModuleList
	recognized := false
	for _, w := range list {
		if w == cmd[2] {
			recognized = true
		}
	}
	if !recognized {
		fmt.Printf("%s[!]%s Error: module does not exist.'\n", tui.RED, tui.RESET)
		return
	}
	// Eventually send command
	s.send(cmd)
	mod := <-s.moduleReqs
	// Switch shell context
	s.Shell.Config.AutoComplete = s.getCompleter("module")
	s.menuContext = "module"
	s.currentModule = mod.ModuleName
	// Add code to change current module in the prompt
}

func (s *Session) reloadModule(cmd []string) {
	s.send(cmd)
	status := <-s.moduleReqs
	fmt.Println()
	fmt.Println(status.Status)
}

func (s *Session) showModuleOptions(cmd []string) {
	s.send(cmd)
	mod := <-s.moduleReqs
	m := mod.Modules[0]

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(3, 50)
	table.SetHeader([]string{"Name", "Value", "Required", "Description"})
	table.SetAutoWrapText(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	table.SetBorder(false)
	// TODO add option for agent alias here
	table.Append([]string{"Agent", tui.Blue(tui.Bold(m.Agent.String())), "true", "Agent on which to run module " + m.Name})
	for _, v := range m.Options {
		table.Append([]string{v.Name, v.Value, strconv.FormatBool(v.Required), v.Description})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) showModuleInfo() {
	s.send(strings.Fields("show options"))
	mod := <-s.moduleReqs
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
	fmt.Println(tui.Dim(wrap(m.Description, 140)))
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
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
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
	fmt.Println(tui.Dim(wrap(m.Notes, 140)))
	fmt.Println()
}

func (s *Session) getModuleList(cmd []string) {
	// Send(cmd)
	mod := <-s.moduleReqs

	list := mod.Modules
	fmt.Println(list)
}

func (s *Session) setModuleOption(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: use 'set <option> <value>'. \n", tui.RED, tui.RESET)
		return
	}
	s.send(cmd)
	opt := <-s.moduleReqs
	if opt.Status != "" {
		fmt.Println()
		fmt.Println(opt.Status)
	}
	if opt.Error != "" {
		fmt.Println()
		fmt.Println(opt.Error)
	}
}

func (s *Session) setAgent(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: Provide an agent to use'. \n", tui.RED, tui.RESET)
		return
	}
	s.send(cmd)
	opt := <-s.moduleReqs
	if opt.Status != "" {
		fmt.Println()
		fmt.Println(opt.Status)
	}
	if opt.Error != "" {
		fmt.Println()
		fmt.Println(opt.Error)
	}
}

func (s *Session) runModule(cmd []string) {
	s.send(cmd)
	mod := <-s.moduleReqs
	fmt.Println(mod)
}

func (s *Session) backModule() {
	s.Shell.Config.AutoComplete = s.getCompleter("main")
	s.menuContext = "main"
	s.currentModule = ""
}
