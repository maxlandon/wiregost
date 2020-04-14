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

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/util"
	"github.com/maxlandon/wiregost/data-service/remote"
)

// Options --------------------------------------------------------------------------
type WorkspaceOptions struct {
	Boundary       []string `long:"boundary" description:"A network range/address, or a list of them, to which scans will be limited"`
	LimitToNetwork bool     `long:"limit-to-network" description:"Limits the scans and other actions to the ranges specified with --boundary"`
	Description    string   `long:"description" description:"A description for this workspace"`
	Default        bool     `long:"default" description:"Make this workspace the default workspace to use on first connection"`
}

// Base --------------------------------------------------------------------------
// Workspace - Root workspace command
type workspaceCmd struct{}

var workspace workspaceCmd

func init() {
	// Command
	MainParser.AddCommand("workspace", "", "", &workspace)
	workspace := MainParser.Find("workspace")
	CommandMap[MAIN_CONTEXT] = append(CommandMap[MAIN_CONTEXT], workspace)
	workspace.ShortDescription = "Manage workspaces (search/add/delete/update)"
	workspace.SubcommandsOptional = true
}

// Execute - Execute base workspace command
func (w *workspaceCmd) Execute(args []string) error {

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
		if w.Name == Context.Workspace.Name {
			name = tui.Bold(tui.Blue(w.Name))
		} else {
			name = w.Name
		}
		data = append(data, []string{name, w.Boundary, def, strconv.FormatBool(w.LimitToNetwork),
			w.UpdatedAt.Format("2006-01-02T15:04:05"), w.Description})
	}

	printWorkspaceTable(data)
	return nil
}

// Add ---------------------------------------------------------------------------
// WorkspaceAddCmd - Add a workspace
type workspaceAddCmd struct {
	Options    *WorkspaceOptions `group:"Workspace properties"`
	Positional struct {
		Name string `description:"Workspace name to add" required:"yes"`
	} `positional-args:"yes" required:"yes"`
}

var workspaceAdd = workspaceAddCmd{Options: &WorkspaceOptions{}}

func init() {
	workspace := MainParser.Find("workspace")

	// Command
	workspace.AddCommand("add", "", "", &workspaceAdd)
	add := workspace.Find("add")
	add.ShortDescription = "Add a workspace, with options (ex: workspace add MyWorkspace --default)"

	// Options
	// boundary := OptionByName("workspace", "add", "boundary")
}

func (wa *workspaceAddCmd) Execute(args []string) error {
	fmt.Println(wa.Options.Description)
	return nil
}

// Switch ---------------------------------------------------------------------------
type workspaceSwitchCmd struct {
	// Arguments
	Positional struct {
		Name string `description:"Workspace to switch to" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var workspaceSwitch workspaceSwitchCmd

func init() {
	workspace := MainParser.Find("workspace")

	workspace.AddCommand("switch", "", "", &workspaceSwitch)
	wsSwitch := workspace.Find("switch")
	wsSwitch.ShortDescription = "Switch to a different workspace"
}

func (ws *workspaceSwitchCmd) Execute(args []string) error {

	workspaces, err := remote.Workspaces(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := range workspaces {
		if workspaces[i].Name == ws.Positional.Name {
			*Context.Workspace = workspaces[i]
			const workspaceID = "workspace_id"
			Context.DBContext = context.WithValue(Context.DBContext, workspaceID, workspaces[i].ID)
			fmt.Printf(Info+"Switched to workspace %s \n", workspaces[i].Name)
		}
	}

	return nil
}

// Delete ---------------------------------------------------------------------------
type workspaceDeleteCmd struct{}

var workspaceDelete workspaceDeleteCmd

func init() {
	workspace := MainParser.Find("workspace")

	workspace.AddCommand("delete", "", "", &workspaceDelete)
	wsDelete := workspace.Find("delete")
	wsDelete.ShortDescription = "Delete a workspace with all its data (hosts, services, credentials, etc.)"
}

func (wd *workspaceDeleteCmd) Execute(args []string) error {
	return nil
}

// Update ---------------------------------------------------------------------------
type workspaceUpdateCmd struct {
	Options    *WorkspaceOptions `group:"Workspace properties"`
	Positional struct {
		Name string `description:"Name of workspace to update"`
	} `positional-args:"yes" required:"yes"`
}

var workspaceUpdate = workspaceUpdateCmd{Options: &WorkspaceOptions{}}

func init() {
	workspace := MainParser.Find("workspace")

	workspace.AddCommand("update", "", "", &workspaceUpdate)
	update := workspace.Find("update")
	update.ShortDescription = "Update a workspace with options (ex: workspace update MyWorkspace --limit-to-network)"
}

func (wu *workspaceUpdateCmd) Execute(args []string) error {
	fmt.Println(OptionByName(*Context.Menu, "workspace", "update", "boundary").Value())
	return nil
}

func printWorkspaceTable(data [][]string) {

	tab := util.NewTable()
	heads := []string{"Name", "Boundary", "Limit", "Default", "UpdatedAt", "Description"}
	widths := []int{10, 20, 5, 5, 10, 40}
	tab.SetColumns(heads, widths)

	tab.AppendBulk(data)
	tab.Output()
}

// import (
//         "context"
//         "fmt"
//         "regexp"
//         "strconv"
//         "strings"
//
//         "github.com/evilsocket/islazy/tui"
//         "github.com/olekukonko/tablewriter"
//
//         "github.com/maxlandon/wiregost/client/util"
//         "github.com/maxlandon/wiregost/data-service/models"
//         "github.com/maxlandon/wiregost/data-service/remote"
// )
//
// func registerWorkspaceCommands() {
//
//         // Declare all commands, subcommands and arguments
//         workspace := &Command{
//                 Name: "workspace",
//                 Help: "Manage workspaces (switch/add/delete/update)",
//                 SubCommands: []*SubCommand{
//                         &SubCommand{Name: "switch", Help: "Change to workspace in which targets and various data are saved"},
//                         &SubCommand{Name: "add", Help: "Add a new workspace (provide a name)"},
//                         &SubCommand{Name: "delete", Help: "Delete a workspace (provide a name)"},
//                         &SubCommand{Name: "update", Help: "Update a workspace (ex: update default --boundary 192.168.1.17)",
//                                 Args: []*Arg{
//                                         &Arg{Name: "name", Type: "string", Required: true},
//                                         &Arg{Name: "limit-to-network", Type: "boolean", Required: false},
//                                         &Arg{Name: "boundary", Type: "string", Required: false},
//                                         &Arg{Name: "description", Type: "string", Required: false},
//                                 },
//                         }},
//                 Handle: func(r *Request) error {
//                         switch length := len(r.Args); {
//                         // No arguments: Print workspaces
//                         case length == 0:
//                                 workspaces(r.context.Workspace)
//                         // Arguments: commands entered
//                         case length >= 1:
//                                 switch r.Args[0] {
//                                 case "switch":
//                                         if len(r.Args) == 2 {
//                                                 switchWorkspace(r.Args[1], r.context.Workspace, &r.context.DBContext, *r.context)
//                                         } else {
//                                                 fmt.Printf("\n" + Error + "Provide a workspace name")
//                                         }
//                                 case "add":
//                                         addWorkspaces(r.Args[1:])
//                                 case "delete":
//                                         deleteWorkspaces(r.Args[1:])
//                                 case "update":
//                                         updateWorkspace(r.Args[1:])
//                                 }
//                         }
//
//                         return nil
//                 },
//         }
//
//         // Add commands for each context
//         AddCommand("main", workspace)
//         AddCommand("module", workspace)
//         AddCommand("ghost", workspace)
//         AddCommand("compiler", workspace)
// }
//
// var WorkspaceOptions struct {
//         LimitToNetwork bool     `long:"limit-to-network" description:"Limit scans and other actions to a network range (--boundary)"`
//         Boundary       []string `long:"boundary" description:"Range of networks/addresses, or a list of both, to use as workspace boundaries"`
//         Description    string   `long:"description" description:"Description of the workspace"`
// }
//
// // func registerWorkspaceCommands() {
// //
// //         // Declare all commands, subcommands and arguments
// //         workspace := &Command{
// //                 Name: "workspace",
// //                 Help: "Manage workspaces (switch/add/delete/update)",
// //                 SubCommands: []*SubCommand{
// //                         &SubCommand{Name: "switch", Help: "Change to workspace in which targets and various data are saved"},
// //                         &SubCommand{Name: "add", Help: "Add a new workspace (provide a name)"},
// //                         &SubCommand{Name: "delete", Help: "Delete a workspace (provide a name)"},
// //                         &SubCommand{Name: "update", Help: "Update a workspace (ex: update default --boundary 192.168.1.17)",
// //                                 Args: []*Arg{
// //                                         &Arg{Name: "name", Type: "string", Required: true},
// //                                         &Arg{Name: "limit-to-network", Type: "boolean", Required: false},
// //                                         &Arg{Name: "boundary", Type: "string", Required: false},
// //                                         &Arg{Name: "description", Type: "string", Required: false},
// //                                 },
// //                         }},
// //                 Handle: func(r *Request) error {
// //                         switch length := len(r.Args); {
// //                         // No arguments: Print workspaces
// //                         case length == 0:
// //                                 workspaces(r.context.Workspace)
// //                         // Arguments: commands entered
// //                         case length >= 1:
// //                                 switch r.Args[0] {
// //                                 case "switch":
// //                                         if len(r.Args) == 2 {
// //                                                 switchWorkspace(r.Args[1], r.context.Workspace, &r.context.DBContext, *r.context)
// //                                         } else {
// //                                                 fmt.Printf("\n" + Error + "Provide a workspace name")
// //                                         }
// //                                 case "add":
// //                                         addWorkspaces(r.Args[1:])
// //                                 case "delete":
// //                                         deleteWorkspaces(r.Args[1:])
// //                                 case "update":
// //                                         updateWorkspace(r.Args[1:])
// //                                 }
// //                         }
// //
// //                         return nil
// //                 },
// //         }
// //
// //         // Add commands for each context
// //         AddCommand("main", workspace)
// //         AddCommand("module", workspace)
// //         AddCommand("ghost", workspace)
// //         AddCommand("compiler", workspace)
// // }
//
// func workspaces(currentWorkspace *models.Workspace) {
//         workspaces, err := remote.Workspaces(nil)
//         if err != nil {
//                 fmt.Println(err.Error())
//         }
//
//         data := [][]string{}
//         for i := range workspaces {
//                 // Default
//                 w := workspaces[i]
//                 def := ""
//                 if w.IsDefault {
//                         def = "default"
//                 }
//                 // Current
//                 name := ""
//                 if w.Name == currentWorkspace.Name {
//                         name = tui.Bold(tui.Blue(w.Name))
//                 } else {
//                         name = w.Name
//                 }
//                 data = append(data, []string{name, w.Description, def, w.Boundary,
//                         strconv.FormatBool(w.LimitToNetwork), w.UpdatedAt.Format("2006-01-02T15:04:05")})
//         }
//
//         printWorkspaceTable(data)
// }
//
// func switchWorkspace(name string, workspace *models.Workspace, ctx *context.Context, sctx ShellContext) {
//         workspaces, err := remote.Workspaces(nil)
//         if err != nil {
//                 fmt.Println(err.Error())
//         }
//         for i := range workspaces {
//                 if workspaces[i].Name == name {
//                         *workspace = workspaces[i]
//                         const workspaceID = "workspace_id"
//                         *ctx = context.WithValue(*ctx, workspaceID, workspaces[i].ID)
//                         workspace = &workspaces[i]
//                         fmt.Printf("\n"+Info+"Switched to workspace %s", workspaces[i].Name)
//                         // Reset currentModule
//                         // *sctx.CurrentModule = ""
//
//                 }
//         }
// }
//
// func addWorkspaces(names []string) {
//         err := remote.AddWorkspaces(nil, names)
//         if err != nil {
//                 fmt.Println(err.Error())
//         } else {
//                 for _, n := range names {
//                         fmt.Printf(Info+"Created workspace %s", n)
//                 }
//         }
// }
//
// func deleteWorkspaces(names []string) {
//         var ids []uint
//         workspaces, _ := remote.Workspaces(nil)
//         for _, w := range workspaces {
//                 for _, name := range names {
//                         if name == w.Name {
//                                 ids = append(ids, w.ID)
//                         }
//                 }
//         }
//
//         err := remote.DeleteWorkspaces(nil, ids)
//         if err != nil {
//                 fmt.Println(err.Error())
//         } else {
//                 for _, n := range names {
//                         fmt.Printf(Info+"Deleted workspace %s", n)
//                 }
//         }
// }
//
// func printWorkspaceTable(data [][]string) {
//         table := util.Table()
//         table.SetColWidth(70)
//         table.SetHeader([]string{"Name", "Description", "Default", "Boundary", "Limit", "Updated At"})
//         table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
//                 tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor})
//         table.SetColMinWidth(1, 40)
//         table.SetColMinWidth(3, 20)
//         table.SetColMinWidth(4, 10)
//         table.AppendBulk(data)
//         table.Render()
// }
//
// func updateWorkspace(args []string) {
//
//         opts := parseOptions(args)
//
//         var w *models.Workspace
//
//         // Check options
//         name, found := opts["workspace_id"]
//         if found {
//                 workspaces, _ := remote.Workspaces(nil)
//                 for i := range workspaces {
//                         if workspaces[i].ID == name.(uint) {
//                                 w = &workspaces[i]
//                         }
//                 }
//         } else {
//                 fmt.Printf(Error + "rovide a workspace name (name='workspace')")
//                 return
//         }
//         desc, found := opts["description"]
//         if found {
//                 w.Description = desc.(string)
//         }
//         boundary, found := opts["boundary"]
//         if found {
//                 w.Boundary = boundary.(string)
//         }
//         limit, found := opts["limit-to-network"]
//         if found {
//                 w.LimitToNetwork = limit.(bool)
//         }
//
//         // Update workspace
//         err := remote.UpdateWorkspace(nil, *w)
//         if err != nil {
//                 fmt.Println(err.Error())
//         } else {
//                 fmt.Printf(Info+"Updated workspace %s", w.Name)
//         }
// }
//
// func parseOptions(args []string) (opts map[string]interface{}) {
//
//         opts = make(map[string]interface{}, 0)
//
//         for _, arg := range args {
//
//                 // LimitToNetwork
//                 if strings.Contains(arg, "limit-to-network") {
//                         vals := strings.Split(arg, "=")
//                         opts["limit-to-network"], _ = strconv.ParseBool(vals[1])
//                 }
//
//                 // Network boundary
//                 if strings.Contains(arg, "boundary") {
//                         vals := strings.Split(arg, "=")
//                         opts["boundary"] = vals[1]
//                 }
//
//                 // Workspace name
//                 if strings.Contains(arg, "name") {
//                         vals := strings.Split(arg, "=")
//                         workspaces, err := remote.Workspaces(nil)
//                         if err != nil {
//                                 fmt.Println(err.Error())
//                         }
//                         for _, w := range workspaces {
//                                 if w.Name == vals[1] {
//                                         opts["workspace_id"] = w.ID
//                                 }
//                         }
//                 }
//                 // Description
//                 if strings.Contains(arg, "description") {
//                         desc := regexp.MustCompile(`\b(description){1}.*"`)
//                         result := desc.FindStringSubmatch(strings.Join(args, " "))
//                         opts["description"] = strings.Trim(strings.TrimPrefix(result[0], "description="), "\"")
//                 }
//         }
//
//         return opts
// }
