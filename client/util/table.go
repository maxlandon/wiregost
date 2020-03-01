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

package util

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

var WrapColumns int

func Table() *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	// Appearance
	table.SetCenterSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetColumnSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)

	// Format
	table.SetAutoWrapText(false)
	table.SetColWidth(30)
	table.SetAutoFormatHeaders(false)
	table.SetReflowDuringAutoWrap(false)

	return table
}

// Function used for description paragraphs and table columns
func Wrap(text string, lineWidth int) (wrapped string) {
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

// SortOptionKeys - Golang prints maps in an ever-changing order, so try at least
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

	return keys
}

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

func SortPostOptions(opts map[string]*clientpb.Option) (keys []string) {

	options := []string{}
	for v, _ := range opts {
		options = append(options, v)
	}

	sort.Strings(options)

	return options
}

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
		tab := Table()
		tab.SetHeader([]string{"Name", "Value", "Required", "Description"})
		tab.SetColWidth(70)
		tab.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		)
		for _, v := range SortListenerOptionKeys(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			description := Wrap(mod.Options[v].Description, WrapColumns-30)
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
		tab = Table()
		tab.SetHeader([]string{"Name", "Value", "Required", "Description"})
		tab.SetColWidth(70)
		tab.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		)
		for _, v := range SortGenerateOptionKeys(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			description := Wrap(mod.Options[v].Description, WrapColumns-30)
			tab.Append([]string{mod.Options[v].Name, mod.Options[v].Value, required, description})
		}
		tab.Render()

	case "post":
		// Print Session
		fmt.Printf(" %sSession: %s%s%s \n", tui.BOLD, "\033[38;5;43m", mod.Options["Session"].Value, tui.RESET)
		fmt.Println()

		tab := Table()
		fmt.Println(tui.Bold(tui.Blue(" Post Options")))
		tab.SetHeader([]string{"Name", "Value", "Required", "Description"})
		tab.SetColWidth(70)
		tab.SetHeaderColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor},
		)
		for _, v := range SortPostOptions(mod.Options) {
			required := ""
			if mod.Options[v].Required {
				required = "yes"
			} else {
				required = "no"
			}
			// Avoid printing session option again
			if v != "Session" {
				description := Wrap(mod.Options[v].Description, WrapColumns-30)
				tab.Append([]string{mod.Options[v].Name, mod.Options[v].Value, required, description})
			}
		}
		tab.Render()
	}
}
