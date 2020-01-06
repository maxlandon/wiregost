package session

import (
	"fmt"
	"os"
	"strconv"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"
)

func (s *Session) useCompiler() {
	// Switch shell context
	s.Shell.Config.AutoComplete = s.getCompleter("compiler")
	s.menuContext = "compiler"
	// Switch prompt
}

func (s *Session) quitCompiler() {
	// Switch prompt
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

func (s *Session) showCompilerOptions(cmd []string) {
	s.send(cmd)
	comp := <-s.compilerReqs

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
	for _, v := range comp.Options {
		table.Append([]string{v.Name, v.Value, strconv.FormatBool(v.Required), v.Description})
	}
	fmt.Println()
	table.Render()
}

func (s *Session) setCompilerOption(cmd []string) {
	s.send(cmd)
	opt := <-s.compilerReqs
	if opt.Status != "" {
		fmt.Println()
		fmt.Println(opt.Status)
	}
	if opt.Error != "" {
		fmt.Println()
		fmt.Println(opt.Error)
	}
}
