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

package console

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	"gopkg.in/yaml.v2"
)

type ConsoleConfig struct {
	Mode        string
	Prompt      string
	HistoryFile string
}

func LoadConsoleConfig() *ConsoleConfig {
	// Load a default console config, eventually parse one if found
	conf := &ConsoleConfig{
		Mode:        "emacs",
		Prompt:      "",
		HistoryFile: "~/.wiregost-client/.history",
	}

	// Load config
	file := filepath.Join(assets.GetRootAppDir(), "console.yaml")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		err = SaveConsoleConfig(conf)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		err = SaveConsoleConfig(conf)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(tui.Red("[!] Error: failed to unmarshal console.yaml file."))
	}

	return conf

}

func SaveConsoleConfig(config *ConsoleConfig) error {
	saveTo := assets.GetRootAppDir()
	configYAML, _ := yaml.Marshal(config)

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write to wiregost-client root directory %s", err))
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "console.yaml"
		saveTo = filepath.Join(saveTo, filename)
	}

	err = ioutil.WriteFile(saveTo, configYAML, 0600)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write config to: %s (%v) \n", saveTo, err))
	}

	f, err := os.OpenFile(saveTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(helpConfig); err != nil {
		panic(err)
	}

	return nil
}

var helpConfig = `

# [ HELP ] ------------------------------------------------------//

# Mode:         vim/emacs
# HistoryFile:  relative/absolute path to history file
# Prompt:       prompt string to use (below are examples and variables)

# Prompt variables:
# {pwd}         Current working directory of the client
# {workspace}   Current Workspace
# {serverip}    IP address of server as specified in client config
# {listeners}   Number of listeners jobs currently running
# {ghosts}      Number of Ghost implants currently connected
# {localip}     IP address of client 

# Prompt colors:
# "{bold}":  BOLD,
# "{dim}":   DIM,
# "{r}":     RED,
# "{g}":     GREEN,
# "{b}":     BLUE,
# "{y}":     YELLOW,
# "{fb}":    FORE BLACK,
# "{fw}":    FORE WHITE,
# "{bdg}":   BACK DARKGRAY,
# "{br}":    BACK RED,
# "{bg}":    BACK GREEN,
# "{by}":    BACK YELLOW,
# "{blb}":   BACK LIGHTBLUE,
# "{reset}": RESET,
# 
# Custom colors:
# "{blink}": BLINKING
# "{lb}":    LIGHTBLUE 
# "{db}":    DARKBLUE 
# "{bddg}":  BACK DARK DARKGRAY 
# "{ly}":    LIGHT YELLOW 

# Examples:

# Emacs mode
# @{lb}{serverip}{reset} in {db}{workspace}{reset} ({jobs},{ghosts})  =>  | @localhost in default (0,2)
#                                                                         |  >  

# Vim mode + module
# @{lb}{serverip}{reset} in {db}{workspace}{reset} ({jobs},{ghosts})  =>  | @localhost in default (0,2) => post(multi)
#                                                                         | [I] >  
`
