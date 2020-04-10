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

package server

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	. "github.com/maxlandon/wiregost/client/commands"
	"github.com/maxlandon/wiregost/client/constants"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/client/transport"
)

// ServerCmd - List available Wiregost servers
type ServerConnectCmd struct {
	Positional struct {
		Server string `description:"Wiregost server to connect to" required:"1"`
	} `positional-args:"yes"`
}

var ServerConnect ServerConnectCmd

func RegisterServerConnect() {
	s := MainParser.Find(constants.Server)
	s.AddCommand(constants.ServerConnect, "", "", &ServerConnect)

	sc := s.Find(constants.ServerConnect)
	sc.ShortDescription = "Connect to one of the available Wiregost servers"
}

// Execute - Command
func (s *ServerConnectCmd) Execute(args []string) error {

	user := strings.Split(s.Positional.Server, "@")[0]
	hostPort := strings.Split(s.Positional.Server, "@")[1]
	lhost := strings.Split(hostPort, ":")[0]
	port := strings.Split(hostPort, ":")[1]
	lport, _ := strconv.Atoi(port)

	configs := assets.GetConfigs()
	var config *assets.ClientConfig
	for _, conf := range configs {
		if (lhost == conf.LHost) && (lport == conf.LPort) && (user == conf.User) {
			config = conf
		}
	}

	fmt.Printf(Warn+"Disconnecting from current server %s:%d \n...\n", Context.Server.Config.LHost, Context.Server.Config.LPort)

	fmt.Printf(Info+"Connecting to %s:%d ...\n", config.LHost, config.LPort)
	send, recv, err := transport.MTLSConnect(config)
	if err != nil {
		errString := fmt.Sprintf(Errorf+"Connection to server failed: %v", err)
		return errors.New(errString)

	}
	fmt.Printf(Success+"Connected to Wiregost server at %s:%d, as user %s%s%s",
		config.LHost, config.LPort, tui.YELLOW, config.User, tui.RESET)
	fmt.Println()

	// Bind connection to server object in console
	Context.Server = core.BindWiregostServer(send, recv)
	Context.Server.Config = config
	go Context.Server.ResponseMapper()

	time.Sleep(time.Millisecond * 50)

	return nil
}
