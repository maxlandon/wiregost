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

package module

import (
	"fmt"
	"sort"
	"strings"

	"github.com/evilsocket/islazy/tui"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/util"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

// ModuleShowOptionsCmd - Show module options
type ModuleShowOptionsCmd struct{}

var ModuleShowOptions ModuleShowOptionsCmd

func RegisterModuleShowOptions() {
	CommandParser.AddCommand(constants.ModuleOptions, "", "", &ModuleShowOptions)

	opts := CommandParser.Find(constants.ModuleOptions)
	CommandMap[MODULE_CONTEXT] = append(CommandMap[MODULE_CONTEXT], opts)
	opts.ShortDescription = "Show options for the current module"
}

// Execute - Run module options
func (so *ModuleShowOptionsCmd) Execute(args []string) error {
	m := Context.Module

	PrintOptions(m)

	return nil
}

// PrintOptions - Print options a for a module, dispatching depending on type
func PrintOptions(mod *clientpb.Module) {

	sub := strings.Join(mod.Path, "/")
	moduleSubtype := ""

	// Get module subtype
	switch subtype := sub; {
	case strings.Contains(subtype, "payload/multi/single"):
		moduleSubtype = "payload/multi/single"
	case strings.Contains(subtype, "payload/multi/stager"):
		moduleSubtype = "payload/multi/stager"
	}

	// Print options depending on module Type/Subtype
	switch mod.Type {
	case "payload":
		// Listener Options
		switch moduleSubtype {
		case "payload/multi/single":
			fmt.Println(tui.Bold(tui.Blue(" Listener Options")))
		case "payload/multi/stager":
			fmt.Println(tui.Bold(tui.Blue(" Staging Listener Options")))
		}

		tab := util.NewTable()
		headers := []string{"Name", "Value", "Required", "Description"}
		widths := []int{15, 15, 7, 40}
		tab.SetColumns(headers, widths)

		tab.SetColWidth(70)

		for _, v := range SortListenerOptionKeys(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			// description := util.Wrap(mod.Options[v].Description, util.WrapColumns-30)
			description := util.AutoWrap(mod.Options[v].Description)
			tab.Append([]string{mod.Options[v].Name, mod.Options[v].Value, required, description})
		}
		tab.Render()

		// Generate Options
		fmt.Println()
		switch moduleSubtype {
		case "payload/multi/single":
			fmt.Println(tui.Bold(tui.Blue(" Implant Options")))
		case "payload/multi/stager":
			fmt.Println(tui.Bold(tui.Blue(" Stager Options")))
		}

		tab = util.NewTable()
		headers = []string{"Name", "Value", "Required", "Description"}
		widths = []int{15, 15, 7, 40}
		tab.SetColumns(headers, widths)
		tab.SetColWidth(70)

		for _, v := range SortGenerateOptionKeys(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			description := util.AutoWrap(mod.Options[v].Description)
			// description := util.Wrap(mod.Options[v].Description, util.WrapColumns-30)
			tab.Append([]string{mod.Options[v].Name, mod.Options[v].Value, required, description})
		}
		tab.Render()

	case "post":
		// Print Session
		fmt.Printf(" %sSession: %s%s%s \n", tui.BOLD, "\033[38;5;43m", mod.Options["Session"].Value, tui.RESET)
		fmt.Println()

		fmt.Println(tui.Bold(tui.Blue(" Post Options")))

		tab := util.NewTable()
		headers := []string{"Name", "Value", "Required", "Description"}
		widths := []int{15, 15, 7, 40}
		tab.SetColumns(headers, widths)

		tab.SetColWidth(70)

		for _, v := range SortPostOptions(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			// Avoid printing session option again
			if v != "Session" {
				description := util.AutoWrap(mod.Options[v].Description)
				// description := util.Wrap(mod.Options[v].Description, util.WrapColumns-30)
				tab.Append([]string{mod.Options[v].Name, mod.Options[v].Value, required, description})
			}
		}
		tab.Render()
	}
}

// SortGenerateOptionKeys - Golang prints maps in an ever-changing order, so try at least
// to give an order for the most important options
func SortGenerateOptionKeys(opts map[string]*clientpb.Option) (keys []string) {

	// Single Payloads
	if _, v := opts["DomainsHTTP"]; v {
		keys = append(keys, "DomainsHTTP")
	}
	if _, v := opts["DomainsMTLS"]; v {
		keys = append(keys, "DomainsMTLS")
	}
	if _, v := opts["DomainsDNS"]; v {
		keys = append(keys, "DomainsDNS")
	}
	if k, v := opts["OS"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Arch"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Format"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["MaxErrors"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["ReconnectInterval"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Save"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["ObfuscateSymbols"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["ListenerDomains"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Canaries"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Debug"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["LimitHostname"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["LimitUsername"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["LimitDatetime"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["LimitDomainJoined"]; v {
		keys = append(keys, k.Name)
	}

	// Stager Payloads
	if _, v := opts["LHostStager"]; v {
		keys = append(keys, "LHostStager")
	}
	if _, v := opts["LPortStager"]; v {
		keys = append(keys, "LPortStager")
	}
	if _, v := opts["StageConfig"]; v {
		keys = append(keys, "StageConfig")
	}
	if _, v := opts["OutputFormat"]; v {
		keys = append(keys, "OutputFormat")
	}
	if _, v := opts["OutputStdout"]; v {
		keys = append(keys, "OutputStdout")
	}
	if _, v := opts["FileName"]; v {
		keys = append(keys, "FileName")
	}
	if _, v := opts["Workspace"]; v {
		keys = append(keys, "Workspace")
	}

	return keys
}

// SortListenerOptionKeys - Listener-specific options
func SortListenerOptionKeys(opts map[string]*clientpb.Option) (keys []string) {

	// Single Payloads
	if _, v := opts["LHost"]; v {
		keys = append(keys, "LHost")
	}
	if _, v := opts["LHost"]; v {
		keys = append(keys, "LPort")
	}
	if _, v := opts["MTLSLHost"]; v {
		keys = append(keys, "MTLSLHost")
	}
	if _, v := opts["MTLSLPort"]; v {
		keys = append(keys, "MTLSLPort")
	}
	if _, v := opts["HTTPLHost"]; v {
		keys = append(keys, "HTTPLHost")
	}
	if _, v := opts["HTTPLPort"]; v {
		keys = append(keys, "HTTPLPort")
	}
	if k, v := opts["Certificate"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["Key"]; v {
		keys = append(keys, k.Name)
	}
	if k, v := opts["LetsEncrypt"]; v {
		keys = append(keys, k.Name)
	}
	if _, v := opts["DomainsHTTPListener"]; v {
		keys = append(keys, "DomainsHTTPListener")
	}
	if k, v := opts["Website"]; v {
		keys = append(keys, k.Name)
	}
	if _, v := opts["DomainsDNSListener"]; v {
		keys = append(keys, "DomainsDNSListener")
	}
	if k, v := opts["DisableCanaries"]; v {
		keys = append(keys, k.Name)
	}
	// Stager Payloads
	if _, v := opts["LHostListener"]; v {
		keys = append(keys, "LHostListener")
		// keys = append(keys, k.Name)
	}
	if _, v := opts["LPortListener"]; v {
		keys = append(keys, "LPortListener")
	}
	if _, v := opts["StageImplant"]; v {
		keys = append(keys, "StageImplant")
	}

	// All
	if k, v := opts["Persist"]; v {
		keys = append(keys, k.Name)
	}

	return keys
}

// // SortPostOptions - Post-module specific options
func SortPostOptions(opts map[string]*clientpb.Option) (keys []string) {

	options := []string{}
	for v := range opts {
		options = append(options, v)
	}

	sort.Strings(options)

	return options
}
