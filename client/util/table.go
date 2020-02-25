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
	"strings"

	"github.com/evilsocket/islazy/tui"
	"github.com/olekukonko/tablewriter"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
)

func Table() *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	// Appearance
	table.SetCenterSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetColumnSeparator(fmt.Sprintf("%s|%s", tui.FOREBLACK, tui.RESET))
	table.SetRowSeparator(tui.Dim("-"))
	table.SetAutoWrapText(true)
	table.SetColWidth(20)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)

	return table
}

// Function used for description paragraphs
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
	if _, v := opts["OutputFormat"]; v {
		keys = append(keys, "OutputFormat")
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
	return keys
}
