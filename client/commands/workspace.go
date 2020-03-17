// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data-service/models"
	"github.com/maxlandon/wiregost/data-service/remote"
)

func registerWorkspaceCommands() {

	// Declare all commands, subcommands and arguments
	workspace := &Command{
		Name: "workspace",
		SubCommands: []string{
			"switch",
			"add",
			"delete",
			"update",
		},
		Args: []*CommandArg{
			&CommandArg{Name: "name", Type: "string", Required: true},
			&CommandArg{Name: "limit-to-network", Type: "boolean", Required: false},
			&CommandArg{Name: "boundary", Type: "string", Required: false},
			&CommandArg{Name: "description", Type: "string", Required: false},
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			// No arguments: Print workspaces
			case length == 0:
				fmt.Println()
				workspaces(r.context.CurrentWorkspace)
				fmt.Println()
			// Arguments: commands entered
			case length >= 1:
				switch r.Args[0] {
				case "switch":
					fmt.Println()
					if len(r.Args) == 2 {
						switchWorkspace(r.Args[1], r.context.CurrentWorkspace, &r.context.Context, *r.context)
					} else {
						fmt.Printf(Error + "Provide a workspace name")
					}
					fmt.Println()
				case "add":
					fmt.Println()
					addWorkspaces(r.Args[1:])
					fmt.Println()
				case "delete":
					fmt.Println()
					deleteWorkspaces(r.Args[1:])
					fmt.Println()
				case "update":
					fmt.Println()
					updateWorkspace(r.Args[1:])
					fmt.Println()
				}
			}

			return nil
		},
	}

	// Add commands for each context
	AddCommand("main", workspace)
	AddCommand("module", workspace)
	AddCommand("ghost", workspace)
	AddCommand("compiler", workspace)
}

func workspaces(currentWorkspace *models.Workspace) {
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	data := [][]string{}
	for i := range workspaces {
		// Default
		w := workspaces[i]
		def := ""
		if w.IsDefault {
			def = "default"
		}
		// Current
		name := ""
		if w.Name == currentWorkspace.Name {
			name = tui.Bold(tui.Blue(w.Name))
		} else {
			name = w.Name
		}
		data = append(data, []string{name, w.Description, def, w.Boundary,
			strconv.FormatBool(w.LimitToNetwork), w.UpdatedAt.Format("2006-01-02T15:04:05")})
	}

	printWorkspaceTable(data)
}

func switchWorkspace(name string, workspace *models.Workspace, ctx *context.Context, sctx ShellContext) {
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := range workspaces {
		if workspaces[i].Name == name {
			*workspace = workspaces[i]
			const workspaceID = "workspace_id"
			*ctx = context.WithValue(*ctx, workspaceID, workspaces[i].ID)
			workspace = &workspaces[i]
			fmt.Printf(Info+"Switched to workspace %s", workspaces[i].Name)
			// Reset currentModule
			*sctx.CurrentModule = ""

		}
	}
}

func addWorkspaces(names []string) {
	err := remote.AddWorkspaces(nil, names)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, n := range names {
			fmt.Printf(Info+"Created workspace %s", n)
		}
	}
}

func deleteWorkspaces(names []string) {
	var ids []uint
	workspaces, _ := remote.Workspaces(nil)
	for _, w := range workspaces {
		for _, name := range names {
			if name == w.Name {
				ids = append(ids, w.ID)
			}
		}
	}

	err := remote.DeleteWorkspaces(nil, ids)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, n := range names {
			fmt.Printf(Info+"Deleted workspace %s", n)
		}
	}
}

func printWorkspaceTable(data [][]string) {
	table := util.Table()
	table.SetColWidth(70)
	table.SetHeader([]string{"Name", "Description", "Default", "Boundary", "Limit", "Updated At"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor})
	table.SetColMinWidth(1, 40)
	table.SetColMinWidth(3, 20)
	table.SetColMinWidth(4, 10)
	table.AppendBulk(data)
	table.Render()
}

func updateWorkspace(args []string) {

	opts := parseOptions(args)

	var w *models.Workspace

	// Check options
	name, found := opts["workspace_id"]
	if found {
		workspaces, _ := remote.Workspaces(nil)
		for i := range workspaces {
			if workspaces[i].ID == name.(uint) {
				w = &workspaces[i]
			}
		}
	} else {
		fmt.Printf(Error + "rovide a workspace name (name='workspace')")
		return
	}
	desc, found := opts["description"]
	if found {
		w.Description = desc.(string)
	}
	boundary, found := opts["boundary"]
	if found {
		w.Boundary = boundary.(string)
	}
	limit, found := opts["limit-to-network"]
	if found {
		w.LimitToNetwork = limit.(bool)
	}

	// Update workspace
	err := remote.UpdateWorkspace(nil, *w)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf(Info+"Updated workspace %s", w.Name)
	}
}

func parseOptions(args []string) (opts map[string]interface{}) {

	opts = make(map[string]interface{}, 0)

	for _, arg := range args {

		// LimitToNetwork
		if strings.Contains(arg, "limit-to-network") {
			vals := strings.Split(arg, "=")
			opts["limit-to-network"], _ = strconv.ParseBool(vals[1])
		}

		// Network boundary
		if strings.Contains(arg, "boundary") {
			vals := strings.Split(arg, "=")
			opts["boundary"] = vals[1]
		}

		// Workspace name
		if strings.Contains(arg, "name") {
			vals := strings.Split(arg, "=")
			workspaces, err := remote.Workspaces(nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			for _, w := range workspaces {
				if w.Name == vals[1] {
					opts["workspace_id"] = w.ID
				}
			}
		}
		// Description
		if strings.Contains(arg, "description") {
			desc := regexp.MustCompile(`\b(description){1}.*"`)
			result := desc.FindStringSubmatch(strings.Join(args, " "))
			opts["description"] = strings.Trim(strings.TrimPrefix(result[0], "description="), "\"")
		}
	}

	return opts
}
