package operator

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

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	// "github.com/bishopfox/sliver/server/transport"
	"github.com/maxlandon/wiregost/cmd"
	"github.com/maxlandon/wiregost/internal/client/console"
	"github.com/maxlandon/wiregost/internal/server/certs"
	"github.com/maxlandon/wiregost/internal/server/configs"
	"github.com/maxlandon/wiregost/internal/server/db"
	"github.com/maxlandon/wiregost/internal/server/db/models"
	"github.com/maxlandon/wiregost/internal/server/multiplayer"
	"github.com/maxlandon/wiregost/internal/server/transport"
)

const (
	// ANSI Colors
	normal    = "\033[0m"
	black     = "\033[30m"
	red       = "\033[31m"
	green     = "\033[32m"
	orange    = "\033[33m"
	blue      = "\033[34m"
	purple    = "\033[35m"
	cyan      = "\033[36m"
	gray      = "\033[37m"
	bold      = "\033[1m"
	clearln   = "\r\x1b[2K"
	upN       = "\033[%dA"
	downN     = "\033[%dB"
	underline = "\033[4m"

	// Info - Display colorful information
	Info = bold + cyan + "[*] " + normal
	// Warn - Warn a user
	Warn = bold + red + "[!] " + normal
	// Debug - Display debug information
	Debug = bold + purple + "[-] " + normal
	// Woot - Display success
	Woot = bold + green + "[$] " + normal
)

// serverOnlyCmds - Server-only operator/multiplayer commands.
func ServerCommands(con *console.Client) (commands []*cobra.Command) {
	startMultiplayer := &cobra.Command{
		Use:     "multiplayer",
		Short:   "Enable multiplayer mode",
		Run:     startMultiplayerModeCmd,
		GroupID: "multiplayer",
	}
	cmd.Bind("multiplayer", false, startMultiplayer, func(f *pflag.FlagSet) {
		f.StringP("lhost", "L", "", "interface to bind server to")
		f.Uint16P("lport", "l", 31337, "tcp listen port")
		f.BoolP("persistent", "p", false, "make persistent across restarts")
	})

	if !con.IsCLI {
		startMultiplayer.GroupID = "core"
	}

	commands = append(commands, startMultiplayer)

	newOperator := &cobra.Command{
		Use:     "new-operator",
		Short:   "Create a new operator config file",
		GroupID: "multiplayer",
		Run:     newOperatorCmd,
	}
	cmd.Bind("operator", false, newOperator, func(f *pflag.FlagSet) {
		f.StringP("lhost", "l", "", "listen host")
		f.Uint16P("lport", "p", 31337, "listen port")
		f.StringP("save", "s", "", "directory/file in which to save config")
		f.StringP("name", "n", "", "operator name")
	})
	cmd.CompleteFlags(newOperator, func(comp *carapace.ActionMap) {
		(*comp)["save"] = carapace.ActionDirectories()
	})

	if !con.IsCLI {
		newOperator.GroupID = "core"
	}

	commands = append(commands, newOperator)

	kickOperator := &cobra.Command{
		Use:     "kick-operator",
		Short:   "Kick an operator from the server",
		GroupID: "multiplayer",
		Run:     kickOperatorCmd,
	}

	cmd.Bind("operator", false, kickOperator, func(f *pflag.FlagSet) {
		f.StringP("name", "n", "", "operator name")
	})

	if !con.IsCLI {
		kickOperator.GroupID = "core"
	}

	commands = append(commands, kickOperator)

	return
}

func newOperatorCmd(cmd *cobra.Command, _ []string) {
	name, _ := cmd.Flags().GetString("name")
	lhost, _ := cmd.Flags().GetString("lhost")
	lport, _ := cmd.Flags().GetUint16("lport")
	save, _ := cmd.Flags().GetString("save")

	if save == "" {
		save, _ = os.Getwd()
	}

	fmt.Printf(Info + "Generating new client certificate, please wait ... \n")
	configJSON, err := multiplayer.NewOperatorConfig(name, lhost, lport)
	if err != nil {
		fmt.Printf(Warn+"%s\n", err)
		return
	}

	saveTo, _ := filepath.Abs(save)
	fi, err := os.Stat(saveTo)
	if !os.IsNotExist(err) && !fi.IsDir() {
		fmt.Printf(Warn+"File already exists %s\n", err)
		return
	}
	if !os.IsNotExist(err) && fi.IsDir() {
		filename := fmt.Sprintf("%s_%s.cfg", filepath.Base(name), filepath.Base(lhost))
		saveTo = filepath.Join(saveTo, filename)
	}
	err = ioutil.WriteFile(saveTo, configJSON, 0o600)
	if err != nil {
		fmt.Printf(Warn+"Failed to write config to: %s (%s) \n", saveTo, err)
		return
	}
	fmt.Printf(Info+"Saved new client config to: %s \n", saveTo)
}

func kickOperatorCmd(cmd *cobra.Command, _ []string) {
	operator, _ := cmd.Flags().GetString("name")

	fmt.Printf(Info+"Removing auth token(s) for %s, please wait ... \n", operator)
	err := db.Session().Where(&models.Operator{
		Name: operator,
	}).Delete(&models.Operator{}).Error
	if err != nil {
		return
	}
	transport.ClearTokenCache()
	fmt.Printf(Info+"Removing client certificate(s) for %s, please wait ... \n", operator)
	err = certs.OperatorClientRemoveCertificate(operator)
	if err != nil {
		fmt.Printf(Warn+"Failed to remove the operator certificate: %v \n", err)
		return
	}
	fmt.Printf(Info+"Operator %s has been kicked out.\n", operator)
}

func startMultiplayerModeCmd(cmd *cobra.Command, _ []string) {
	lhost, _ := cmd.Flags().GetString("lhost")
	lport, _ := cmd.Flags().GetUint16("lport")
	persistent, _ := cmd.Flags().GetBool("persistent")

	_, err := multiplayer.StartClientListener(lhost, lport)
	if err == nil {
		fmt.Printf(Info + "Multiplayer mode enabled!\n")
		if persistent {
			serverConfig := configs.GetServerConfig()
			serverConfig.AddMultiplayerJob(&configs.MultiplayerJobConfig{
				Host: lhost,
				Port: lport,
			})
			serverConfig.Save()
		}
	} else {
		fmt.Printf(Warn+"Failed to start job %v\n", err)
	}
}
