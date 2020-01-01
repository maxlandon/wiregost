package session

import (
	"fmt"
	"os"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) StackShow() {
	s.Send(strings.Fields("stack show"))
	stack := <-s.moduleReqs

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

func (s *Session) StackPop(cmd []string) {
	s.Send(cmd)
	switch len(cmd) {
	case 2:
		s.currentModule = ""
	case 3:
		if strings.ToLower(cmd[2]) == strings.ToLower(s.currentModule) {
			s.currentModule = ""
		}
	}
	// Temporary: return to main menu completion.
	// This will change when the code will handle fallback on next module in stack.
	s.Shell.Config.AutoComplete = s.getCompleter("main")

}
