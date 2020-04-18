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

// Console - Central object of the client UI
type Console struct{}

// NewConsole - Instantiates a console with some default behavior
func NewConsole() *Console {

	return &Console{}
}

// Setup - Setup various elements of the console.
func (c *Console) Setup() {

}

// ShareContext - The console exposes its context to other packages
func (c *Console) ShareContext() {

}

// Start - Start the console
func (c *Console) Start() {

}

// Refresh - Computes prompt and current context
func (c *Console) Refresh() {

}

// Exit - Kill the current client console
func (c *Console) Exit() {

}
