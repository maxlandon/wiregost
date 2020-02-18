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

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/core"
	"github.com/maxlandon/wiregost/client/transport"
	. "github.com/maxlandon/wiregost/client/util"
)

func getDefaultServerConfig() *assets.ClientConfig {
	configs := assets.GetConfigs()
	if len(configs) == 0 {
		fmt.Printf(Warnf+"No config files found at %s or -config\n", assets.GetConfigDir())
		return nil
	}

	var config *assets.ClientConfig
	for _, conf := range configs {
		if conf.IsDefault {
			config = conf
		}
	}

	return config
}

func (c *Console) connect(config *assets.ClientConfig) error {

	// Initiate connection
	fmt.Printf(Info+"Connecting to %s:%d ...\n", config.LHost, config.LPort)
	send, recv, err := transport.MTLSConnect(config)
	if err != nil {
		errString := fmt.Sprintf(Errorf+"Connection to server failed: %v", err)
		return errors.New(errString)
	} else {
		fmt.Printf(Success+"Connected to Wiregost server at %s:%d, as user %s%s%s",
			config.LHost, config.LPort, tui.YELLOW, config.User, tui.RESET)
		fmt.Println()
	}

	// Bind connection to server object in console
	c.server = core.BindWiregostServer(send, recv)
	go c.server.ResponseMapper()

	// Actualize shell context with server
	c.shellContext.Server = c.server
	c.shellContext.Server.Config = config

	return nil
}
