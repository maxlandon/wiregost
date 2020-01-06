package session

import (
	"fmt"
	"os"
	"strings"

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

func (s *Session) stackPop(cmd []string) {
	s.send(cmd)
	// Wait for new current module fallback
	fallback := <-s.moduleReqs
	if fallback.ModuleName != "" {
		if s.currentModule != "" {
			s.currentModule = fallback.ModuleName
		}
	} else {
		s.currentModule = ""
		s.Shell.Config.AutoComplete = s.getCompleter("main")
	}

}

func (s *Session) stackUse(cmd []string) {
	s.send([]string{"use", "module", cmd[2]})
	mod := <-s.moduleReqs
	// Switch shell context
	s.Shell.Config.AutoComplete = s.getCompleter("module")
	s.menuContext = "module"
	s.currentModule = mod.ModuleName
}
