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

package models

import "github.com/google/uuid"

type Module struct {
	Agent        uuid.UUID   // The Agent that will later be associated with this module prior to execution
	Name         string      `json:"name"`                 // Name of the module
	Type         string      `json:"type"`                 // Type of module (i.e. standard or extended)
	Author       []string    `json:"author"`               // A list of module authors
	Credits      []string    `json:"credits"`              // A list of people to credit for underlying tool or techniques
	Path         []string    `json:"path"`                 // Path to the module (i.e. data/modules/powershell/powerview)
	Platform     string      `json:"platform"`             // Platform the module can run on (i.e. Windows, Linux, Darwin, or ALL)
	Arch         string      `json:"arch"`                 // The Architecture the module can run on (i.e. x86, x64, MIPS, ARM, or ALL)
	Lang         string      `json:"lang"`                 // What language does the module execute in (i.e. PowerShell, Python, or Perl)
	Priv         bool        `json:"privilege"`            // Does this module required a privileged level account like root or SYSTEM?
	Description  string      `json:"description"`          // A description of what the module does
	Notes        string      `json:"notes"`                // Additional information or notes about the module
	Commands     []string    `json:"commands"`             // A list of commands to be run on the agent
	SourceRemote string      `json:"remote"`               // Online or remote source code for a module (i.e. https://raw.githubusercontent.com/PowerShellMafia/PowerSploit/master/Exfiltration/Invoke-Mimikatz.ps1)
	SourceLocal  []string    `json:"local"`                // The local file path to the script or payload
	Options      []Option    `json:"options"`              // A list of configurable options/arguments for the module
	Powershell   interface{} `json:"powershell,omitempty"` // An option json object containing commands and configuration items specific to PowerShell
}

// Option is a structure containing the keys for the object
type Option struct {
	Name        string `json:"name"`        // Name of the option
	Value       string `json:"value"`       // Value of the option
	Required    bool   `json:"required"`    // Is this a required option?
	Flag        string `json:"flag"`        // The command line flag used for the option
	Description string `json:"description"` // A description of the option
}
