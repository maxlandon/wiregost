package cli

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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/maxlandon/wiregost/internal/server/certs"
	"github.com/maxlandon/wiregost/internal/server/multiplayer"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Generate operator configuration files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString(nameFlagStr)
		if err != nil {
			fmt.Printf("Failed to parse --%s flag %s", nameFlagStr, err)
			os.Exit(1)
		}
		if name == "" {
			fmt.Printf("Must specify --%s", nameFlagStr)
			os.Exit(1)
		}

		lhost, err := cmd.Flags().GetString(lhostFlagStr)
		if err != nil {
			fmt.Printf("Failed to parse --%s flag %s", lhostFlagStr, err)
			os.Exit(1)
		}
		if lhost == "" {
			fmt.Printf("Must specify --%s", lhostFlagStr)
			os.Exit(1)
		}

		lport, err := cmd.Flags().GetUint16(lportFlagStr)
		if err != nil {
			fmt.Printf("Failed to parse --%s flag %s", lportFlagStr, err)
			os.Exit(1)
		}

		save, err := cmd.Flags().GetString(saveFlagStr)
		if err != nil {
			fmt.Printf("Failed to parse --%s flag %s", saveFlagStr, err)
			os.Exit(1)
		}
		if save == "" {
			save, _ = os.Getwd()
		}

		certs.SetupCAs()
		configJSON, err := multiplayer.NewOperatorConfig(name, lhost, lport)
		if err != nil {
			fmt.Printf("Failed: %s\n", err)
			os.Exit(1)
		}

		saveTo, _ := filepath.Abs(save)
		fi, err := os.Stat(saveTo)
		if !os.IsNotExist(err) && !fi.IsDir() {
			fmt.Printf("File already exists: %s\n", err)
			os.Exit(1)
		}
		if !os.IsNotExist(err) && fi.IsDir() {
			filename := fmt.Sprintf("%s_%s.cfg", filepath.Base(name), filepath.Base(lhost))
			saveTo = filepath.Join(saveTo, filename)
		}
		err = os.WriteFile(saveTo, configJSON, 0o600)
		if err != nil {
			fmt.Printf("Write failed: %s (%s)\n", saveTo, err)
			os.Exit(1)
		}
	},
}
