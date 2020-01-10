package session

import (
	// Standard
	"fmt"
	"os"
	"strings"

	// 3rd party
	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) stackShow() {
	s.send(strings.Fields("stack show"))
	stack := <-s.moduleReqs

	// Print stack
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator(tui.Dim("|"))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetColumnSeparator(tui.Dim("|"))
	table.SetColMinWidth(1, 50)
	table.SetHeader([]string{"Name", "Source Local", "Language"})
	table.SetAutoWrapText(true)
	table.SetReflowDuringAutoWrap(true)
	table.SetColWidth(80)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
	)

	table.SetBorder(false)
	// TODO add option for agent alias here
	for i := len(stack.Modules) - 1; i >= 0; i-- {
		if strings.ToLower(strings.TrimSuffix(strings.Join(stack.Modules[i].Path, "/"), ".json")) == strings.ToLower(s.currentModule) {
			table.Rich([]string{stack.Modules[i].Name, strings.TrimPrefix(strings.Join(stack.Modules[i].SourceLocal, "/"), "data/src"), stack.Modules[i].Lang},
				[]tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor},
				})
		} else {
			table.Append([]string{stack.Modules[i].Name, strings.TrimPrefix(strings.Join(stack.Modules[i].SourceLocal, "/"), "data/src"), stack.Modules[i].Lang})
		}
	}
	table.Render()
}

func (s *Session) stackPop(cmd []string) {
	var popped string
	if len(cmd) == 3 {
		// Check module, so we don't send garbage to the server.
		s.send([]string{"stack", "list"})
		resp := <-s.moduleReqs
		list := resp.ModuleList
		recognized := false
		for _, m := range list {
			if m == cmd[2] {
				recognized = true
				popped = cmd[2]
			}
		}
		if !recognized {
			fmt.Printf("%s[!]%s Error in module name: not in stack.'\n", tui.RED, tui.RESET)
			return
		}
	} else {
		popped = s.currentModule
	}
	// Eventually send command
	s.send(cmd)
	// Wait for new current module fallback
	fallback := <-s.moduleReqs
	if fallback.ModuleName != "" {
		if s.currentModule != "" && s.currentModule == popped {
			s.currentModule = fallback.ModuleName
		}
	} else {
		s.currentModule = ""
		s.Shell.Config.AutoComplete = s.getCompleter("main")
	}

}

func (s *Session) stackUse(cmd []string) {
	if len(cmd) < 3 {
		fmt.Printf("%s[!]%s Invalid command: give a stack module to use'\n", tui.RED, tui.RESET)
		return
	}
	// Check module, so we don't send garbage to the server.
	s.send([]string{"stack", "list"})
	resp := <-s.moduleReqs
	list := resp.ModuleList
	recognized := false
	for _, m := range list {
		if m == cmd[2] {
			recognized = true
		}
	}
	if !recognized {
		fmt.Printf("%s[!]%s Error in module name: not in stack.'\n", tui.RED, tui.RESET)
		return
	}
	s.send([]string{"use", "module", cmd[2]})
	mod := <-s.moduleReqs
	// Switch shell context
	s.Shell.Config.AutoComplete = s.getCompleter("module")
	s.menuContext = "module"
	s.currentModule = mod.ModuleName
}
