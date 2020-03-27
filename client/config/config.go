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

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	"gopkg.in/yaml.v2"
)

// Config - Stores various configuration elements for the shell
type Config struct {
	Mode                  string
	Prompt                string
	Wrap                  string
	ImplantPrompt         string
	HistoryFile           string
	SessionPathCompletion bool
}

// LoadConsoleConfig - Loads the config file for the console and parse it
func LoadConsoleConfig() *Config {
	// Load a default console config, eventually parse one if found
	conf := &Config{
		Mode:                  "emacs",
		Prompt:                "",
		Wrap:                  "large",
		ImplantPrompt:         "",
		HistoryFile:           "~/.wiregost-client/.history",
		SessionPathCompletion: false,
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

// SaveConsoleConfig - Save the config to file
func SaveConsoleConfig(config *Config) error {
	saveTo := assets.GetRootAppDir()
	configYAML, _ := yaml.Marshal(config)

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Cannot write to wiregost-client root directory %s", err)
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "console.yaml"
		saveTo = filepath.Join(saveTo, filename)
	}

	err = ioutil.WriteFile(saveTo, configYAML, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write config to: %s (%v) \n", saveTo, err)
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

# Mode:                     vim/emacs - Input mode.
# HistoryFile:              relative/absolute path to history file.
# Prompt:                   prompt string to use (below are examples and variables).
# Wrap:                     small/large - overall columns width, to accomodate for little screens
#                           or for those who may need multiple consoles on the same screen.
# ImplantPrompt:            prompt string to use when interacting with an implant.
# SessionPathCompletion:    Enable session's path completer (will send ListDir requests often).

# Prompt variables:
# {pwd}         Current working directory of the client
# {workspace}   Current Workspace
# {serverip}    IP address of server as specified in client config
# {listeners}   Number of listeners jobs currently running
# {ghosts}      Number of Ghost implants currently connected
# {localip}     IP address of client 

# ImplantPrompt variables:
# {user}        Username of the implant process owner
# {host}        Hostname of the target 
# {rpwd}        Current working directory of the implant 
# {os}          Operating System of the target 
# {arch}        CPU architecture of the target 
# {uid}         User ID of the implant process owner 
# {gid}         Group ID of the implant process owner 
# {pid}         Process ID of the implant 
# {transport}   C2 protocol used by the implant (mtls/https/dns) 
# {address}     host:port address of the implant 

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

# Notes:
# When setting the implant prompt, you can also mix Prompt and ImplantPrompt variables as you wish.

# Examples:

# Emacs mode
# @{lb}{serverip}{reset} in {db}{workspace}{reset} ({jobs},{ghosts})  =>  | @localhost in default (0,2)
#                                                                         |  >  

# Vim mode + module
# @{lb}{serverip}{reset} in {db}{workspace}{reset} ({jobs},{ghosts})  =>  | @localhost in default (0,2) => post(multi)
#                                                                         | [I] >  
`
