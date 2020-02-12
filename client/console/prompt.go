// Wiregost - Post-Exploitation & Implant Framework
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

package console

var (
	// Prompt - The prompt object used by the console
	Prompt prompt
)

// prompt - Stores all variables necessary to the console prompt
type prompt struct {
	// Strings
	Base struct { // Non custom prompt
		Main  string
		Ghost string
	}
	Custom struct { // Custom prompt
		Main  string
		Ghost string
	}
	Multiline struct { // Multiline prompts
		Vim   string
		Emacs string
	}
	// Variables
	Workspace *string
	Module    *string
	Menu      *string
	ServerIP  *string
	// Callbacks & Colors
	Callbacks map[string]func() string
	Effects   map[string]string
}

// SetPrompt - Initializes the Prompt object
func (c *console) SetPrompt() {

	Prompt = prompt{} // Initialize
	setCallbacks(Prompt)

	return
}

// setCallbacks - Initializes all callbacks for prompt
func setCallbacks(prompt prompt) {

}

// render - Computes all variables and outputs prompt
func (p *prompt) render() (prompt string, multi string) {
	return
}
