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
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/help"
	"github.com/maxlandon/wiregost/client/util"
)

func RegisterCoreCommands() {

	// Shell -----------------------------------------------------------------------------------------//

	shell := &Command{
		Name: "!",
		Help: help.GetHelpFor("!"),
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				fmt.Printf(help.GetHelpFor("!"))
				fmt.Println()
			default:
				fmt.Println()
				out, err := util.Shell(r.Args[0:])
				if err != nil {
					fmt.Printf("%s Error:%s %s", tui.RED, tui.RESET, err)
				} else {
					fmt.Printf("%s%s\n", tui.RESET, out)
				}
			}
			return nil
		},
	}

	// Add shell command
	AddCommand("main", shell)
	AddCommand("module", shell)
	AddCommand("ghost", shell)
	AddCommand("compiler", shell)

	// Change directory ------------------------------------------------------------------------------//

	cd := &Command{
		Name: "cd",
		Help: help.GetHelpFor("cd"),
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				fmt.Printf(help.GetHelpFor("cd"))
				fmt.Println()
			default:
				fmt.Println()
				dir, err := fs.Expand(r.Args[0])
				err = os.Chdir(dir)
				if err != nil {
					fmt.Printf("%s Error:%s %s", tui.RED, tui.RESET, err)
				} else {
					fmt.Printf("%s*%s Changed directory to %s", tui.BLUE, tui.RESET, dir)
					fmt.Println()
				}
			}
			return nil
		},
	}

	// Add cd command
	AddCommand("main", cd)
	AddCommand("module", cd)
	AddCommand("ghost", cd)
	AddCommand("compiler", cd)

	// Resource file ------------------------------------------------------------------------------//

	resource := &Command{
		Name: "resource",
		Help: help.GetHelpFor("resource"),
		SubCommands: []string{
			"make",
			"load",
		},
		Args: []*CommandArg{
			&CommandArg{Name: "filename", Type: "string"},
			&CommandArg{Name: "length", Type: "int"},
		},
		Handle: func(r *Request) error {
			switch length := len(r.Args); {
			case length == 0:
				fmt.Println()
				fmt.Printf(help.GetHelpFor("resource"))

			// Arguments: commands entered
			case length >= 1:
				fmt.Println()
				switch r.Args[0] {
				// Make a resource file --------------------------------------------------------------------------------------//
				case "make":
					if len(r.Args) < 3 {
						fmt.Printf("%s[!]%s Missing some parameters (type 'resource' for help)", tui.YELLOW, tui.RESET)
						fmt.Println()
						return nil
					}
					var filename string
					var length int
					for _, arg := range r.Args[1:] {
						if strings.Contains(arg, "filename") {
							filename = strings.Split(arg, "=")[1]
						}
						if !strings.Contains(filename, "rc") {
							filename = filename + ".rc"
						}
						if strings.Contains(arg, "length") {
							l := strings.Split(arg, "=")[1]
							length, _ = strconv.Atoi(l)
						}
					}
					if filename == "" {
						fmt.Printf("%s[!]%s Missing resource filename (filename='name.rc')", tui.YELLOW, tui.RESET)
						return nil
					}
					if length == 0 {
						fmt.Printf("%s[!]%s Missing resource command length (length=8)", tui.YELLOW, tui.RESET)
						return nil
					}
					if !strings.Contains(filename, "rc") {
						filename = filename + ".rc"
					}

					// Check if resource already exists
					file, _ := fs.Expand(path.Join(assets.GetResourceDir(), filename))
					if fs.Exists(file) {
						fmt.Printf("%sError:%s resource file already exists. Cannot overwrite it.",
							tui.RED, tui.RESET)
						return nil
					}
					// If not, create it
					res, err := os.Create(file)
					if err != nil {
						panic(err)
					}
					defer res.Close()

					// Find file length, for subsequent selection
					hist, _ := os.Open(assets.GetRootAppDir() + "/.history")
					hlength := 0
					scanner := bufio.NewScanner(hist)
					for scanner.Scan() {
						hlength += 1
					}
					hist.Close()

					// Select last x commands in file and output them to resource file
					hist, _ = os.Open(assets.GetRootAppDir() + "/.history")
					defer hist.Close()
					scan := bufio.NewScanner(hist)
					for scan.Scan() {
						if hlength <= length {
							res.Write([]byte(scan.Text() + "\n"))
						}
						hlength -= 1
					}
					fmt.Printf("%s[*]%s Resource file created and filed with last %s commands.%s",
						tui.GREEN, tui.RESET, strconv.Itoa(length), tui.RESET)

					// Load a resource file --------------------------------------------------------------------------------------//
				case "load":
					if len(r.Args) == 1 {
						fmt.Printf("%s[!]%s Missing resource filename (resource load <file>)", tui.YELLOW, tui.RESET)
						return nil
					}
					filename := r.Args[1]
					filestr, _ := fs.Expand(path.Join(assets.GetResourceDir(), filename))
					file, _ := os.Open(filestr)
					if filepath.Ext(filestr) != ".rc" {
						fmt.Printf("%s[!]%s File must be a configuration (.rc) file.", tui.RED, tui.RESET)
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						lign := scanner.Text()
						cmds := strings.Split(lign, " ")
						command := FindCommand("main", cmds[0])
						if command != nil {
							command.Handle(NewRequest(command, cmds[1:], r.context))
							fmt.Println(tui.Dim(tui.Blue("------------------------------------------------------")))
						} else {
							fmt.Println()
							fmt.Printf("%sError:%s %s%s%s is not a valid command.",
								tui.RED, tui.RESET, tui.YELLOW, cmds[0], tui.RESET)
						}
					}
					if err := scanner.Err(); err != nil {
						fmt.Printf("%s[!]%s Error parsing resource command %s%s%s : %s",
							tui.RED, tui.RESET, tui.YELLOW, scanner.Text(), tui.RESET, err)
					}
				}
			}
			fmt.Println()
			return nil
		},
	}

	// Add resource command
	AddCommand("main", resource)
	AddCommand("module", resource)
	AddCommand("ghost", resource)
}
