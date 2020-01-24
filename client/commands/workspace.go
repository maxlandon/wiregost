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
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	"github.com/maxlandon/wiregost/client/completers"
	consts "github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data_service/models"
	"github.com/maxlandon/wiregost/data_service/remote"
)

type currentWorkspace struct {
	Workspace *models.Workspace
	observers []observer
}

var CurrentWorkspace = &currentWorkspace{
	observers: []observer{},
}

func RegisterWorkspaceCommands(workspace *models.Workspace, cctx *context.Context, app *grumble.App) {

	// Base command, list workspaces
	workspaceCommand := &grumble.Command{
		Name:     consts.WorkspaceStr,
		Help:     tui.Dim("Manage Wiregost workspaces"),
		LongHelp: help.GetHelpFor(consts.WorkspaceStr),
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			workspaces(gctx)
			fmt.Println()
			return nil
		},
		HelpGroup: consts.DataServiceHelpGroup,
	}

	// Switch workspace
	workspaceCommand.AddCommand(&grumble.Command{
		Name:      "switch",
		Help:      tui.Dim("Switch to workspace"),
		LongHelp:  help.GetHelpFor(consts.WorkspaceStr),
		AllowArgs: true,
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			switchWorkspace(gctx.Args[0], workspace, cctx, gctx)
			fmt.Println()
			return nil
		},
		Completer: completers.CompleteWorkspaces(),
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Add workspaces
	workspaceCommand.AddCommand(&grumble.Command{
		Name:      "add",
		Help:      tui.Dim("Add one or more workspaces"),
		LongHelp:  help.GetHelpFor(consts.WorkspaceStr),
		AllowArgs: true,
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			addWorkspaces(gctx)
			fmt.Println()
			return nil
		},
		Flags: func(f *grumble.Flags) {
			f.StringL("name", "", "Name of workspace to add")
			f.StringL("description", "", "A description for the workspace")
			f.StringL("boundary", "", "One or several IPv4/IPv6 Addresses/Ranges, comma-separated (ex: 192.168.1.15,230.16.13.15)")
			f.BoolL("limit_to_network", false, "Limit tools activities (exploits, scanners, etc) to the workspace's boundary")
		},
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Delete workspaces
	workspaceCommand.AddCommand(&grumble.Command{
		Name:      "delete",
		Help:      tui.Dim("Delete one or more workspaces"),
		LongHelp:  help.GetHelpFor(consts.WorkspaceStr),
		AllowArgs: true,
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			deleteWorkspaces(gctx)
			fmt.Println()
			return nil
		},
		Completer: completers.CompleteWorkspaces(),
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Update workspace
	workspaceCommand.AddCommand(&grumble.Command{
		Name:      "update",
		Help:      tui.Dim("Update a workspace with provided options"),
		LongHelp:  help.GetHelpFor(consts.WorkspaceStr),
		AllowArgs: true,
		Run: func(gctx *grumble.Context) error {
			fmt.Println()
			updateWorkspace(gctx)
			fmt.Println()
			return nil
		},
		Flags: func(f *grumble.Flags) {
			f.StringL("name", "", "Workspace to update")
			f.StringL("description", "", "A description for the workspace")
			f.StringL("boundary", "", "One or several IPv4/IPv6 Addresses/Ranges, comma-separated (ex: 192.168.1.15,230.16.13.15)")
			f.BoolL("limit_to_network", false, "Limit tools activities (exploits, scanners, etc) to the workspace's boundary")
		},
		// Completer: completers.CompleteWorkspacesAndFlags(),
		HelpGroup: consts.DataServiceHelpGroup,
	})

	// Register root workspace command
	app.AddCommand(workspaceCommand)
}

func workspaces(gctx *grumble.Context) {
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	data := [][]string{}
	for i, _ := range workspaces {
		// Default
		w := workspaces[i]
		def := ""
		if w.IsDefault {
			def = "default"
		}
		// Current
		name := ""
		if w.Name == CurrentWorkspace.Workspace.Name {
			name = tui.Bold(tui.Blue(w.Name))
		} else {
			name = w.Name
		}
		data = append(data, []string{name, w.Description, def, w.Boundary,
			strconv.FormatBool(w.LimitToNetwork), w.UpdatedAt.Format("2006-01-02T15:04:05")})
	}

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

func switchWorkspace(name string, workspace *models.Workspace, cctx *context.Context, gctx *grumble.Context) {
	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, _ := range workspaces {
		if workspaces[i].Name == name {
			*workspace = workspaces[i]
			*cctx = context.WithValue(*cctx, "workspace_id", workspaces[i].ID)
			CurrentWorkspace.SetCurrentWorkspace(workspaces[i])
			fmt.Printf("%s*%s Switched to workspace %s\n",
				tui.BLUE, tui.RESET, workspaces[i].Name)
		}
	}
}

func addWorkspaces(gctx *grumble.Context) {
	names := gctx.Args
	err := remote.AddWorkspaces(nil, names)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, n := range names {
			fmt.Printf("%s*%s Created workspace %s\n",
				tui.BLUE, tui.RESET, n)
		}
	}
}

func deleteWorkspaces(gctx *grumble.Context) {
	names := gctx.Args
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
			fmt.Printf("%s*%s Deleted workspace %s\n",
				tui.BLUE, tui.RESET, n)
		}
	}
}

func updateWorkspace(gctx *grumble.Context) {

	var workspace *models.Workspace

	if gctx.Flags.String("name") != "" {
		workspaces, _ := remote.Workspaces(nil)
		for i, _ := range workspaces {
			if workspaces[i].Name == gctx.Flags.String("name") {
				workspace = &workspaces[i]
			}
		}
	} else {
		fmt.Printf("%s[!]%s Provide a workspace name (--name 'workspace')\n",
			tui.RED, tui.RESET)
		return
	}

	if gctx.Flags.String("description") != "" {
		workspace.Description = gctx.Flags.String("description")
	}
	if gctx.Flags.String("boundary") != "" {
		workspace.Boundary = gctx.Flags.String("boundary")
	}
	if gctx.Flags.Bool("limit_to_network") {
		workspace.LimitToNetwork = gctx.Flags.Bool("limit_to_network")
	}

	err := remote.UpdateWorkspace(nil, *workspace)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%s*%s Update workspace %s\n",
			tui.BLUE, tui.RESET, workspace.Name)
	}
}

func (cw *currentWorkspace) AddObserver(fn observer) {
	cw.observers = append(cw.observers, fn)
}

func (cw *currentWorkspace) SetCurrentWorkspace(workspace models.Workspace) {
	cw.Workspace = &workspace
	for _, fn := range cw.observers {
		fn()
	}
}
