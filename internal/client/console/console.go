package console

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
	"context"
	"fmt"
	"log"
	"os"

	"github.com/reeflective/console"
	"golang.org/x/exp/slog"

	"github.com/maxlandon/wiregost/internal/client/assets"
	"github.com/maxlandon/wiregost/internal/client/version"
	"github.com/maxlandon/wiregost/internal/proto/commonpb"
	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
)

// Client is the client console application for wiregost.
type Client struct {
	App      *console.Console
	Rpc      rpcpb.CoreClient
	Settings *assets.ClientSettings
	IsServer bool
	IsCLI    bool

	printf        func(format string, args ...any) (int, error)
	jsonHandler   slog.Handler
	logFile       *os.File
	asciicastFile *os.File
}

// NewClient sets up a new console application without starting it.
func NewClient(isServer bool) *Client {
	assets.Setup(false, false)
	settings, _ := assets.LoadSettings()

	con := &Client{
		App:      console.New("wiregost"),
		Settings: settings,
		IsServer: isServer,
	}

	con.App.NewlineBefore = true
	con.App.NewlineAfter = true

	return con
}

// Setup binds a working RPC connection to the client, and optionally some commands.
// It makes use of the client RPC configuration for banner information display.
func Setup(con *Client, cfg *assets.ClientConfig, rpc rpcpb.CoreClient, commands console.Commands) {
	con.Rpc = rpc

	// The console application needs to query the terminal for cursor positions
	// when asynchronously printing logs (that is, when no command is running).
	// If ran from a system shell, however, those queries will block because
	// the system shell is in control of stdin. So just use the classic Printf.
	if con.IsCLI {
		con.printf = fmt.Printf
	} else {
		con.printf = con.App.TransientPrintf
	}

	con.App.SetPrintLogo(func(_ *console.Console) {
		con.printBanner(cfg)
	})

	// console logger
	if con.Settings.ConsoleLogs {
		// Classic logs
		con.logFile = getConsoleLogFile()
		consoleLogStream, err := con.ClientLogStream("json")
		if err != nil {
			log.Printf("Could not get client log stream: %s", err)
		}
		con.setupLogger(con.logFile, consoleLogStream)

		// Ascii cast sessions (complete terminal interface).
		con.asciicastFile = getConsoleAsciicastFile()

		asciicastStream, err := con.ClientLogStream("asciicast")
		con.setupAsciicastRecord(con.asciicastFile, asciicastStream)
	}

	// Bind commands
	con.App.Menu("").SetCommands(commands)
}

// Start starts the console application (blocking).
func (c *Client) Start() error {
	// Close log files on exit
	defer c.logFile.Close()
	defer c.asciicastFile.Close()

	return c.App.Start()
}

// printBanner - Print Wiregost banner, client & server information, etc.
func (c *Client) printBanner(cfg *assets.ClientConfig) {
	// User info
	var username string
	if c.IsServer {
		username = "server (no user)"
	} else {
		username = cfg.Operator
	}

	userStr := fmt.Sprintf("Connected as user: %s%s%s", Orange, username, Normal)
	user := fmt.Sprintf("%-91s", userStr)

	// Server version
	var lhost string
	if c.IsServer {
		lhost = "local"
	} else {
		lhost = fmt.Sprintf("%s:%d", cfg.LHost, cfg.LPort)
	}

	serverStr := fmt.Sprintf("Server connection: %sok%s (%s%s%s)", Green, Normal, Blue, lhost, Normal)
	server := fmt.Sprintf("%-100s", serverStr)

	serverVer, err := c.Rpc.GetVersion(context.Background(), &commonpb.Empty{})
	if err != nil {
		panic(err.Error())
	}

	serverSemVer := fmt.Sprintf("%d.%d.%d", serverVer.Major, serverVer.Minor, serverVer.Patch)
	serVerStr := fmt.Sprintf("Server version: %s%s%s [%s%s%s]", Orange, serverSemVer, Normal, LightRed, serverVer.GetCommit(), Normal)
	serVer := fmt.Sprintf("%-106s", serVerStr)

	var cliVer string

	// Client version
	if c.IsServer {
		cliVer = serVer
	} else {
		v := version.SemanticVersion()
		cliSemVer := fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
		cliVerStr := fmt.Sprintf("Client version: %s%s%s [%s%s%s]", Orange, cliSemVer, Normal, LightRed, version.GitCommit, Normal)
		cliVer = fmt.Sprintf("%-106s", cliVerStr)
	}

	// Setup banners
	bannerServerConnection := fmt.Sprintf(`                                                                                                            ,,.,=++============+,               
                                                                                                       ,+==========================..           
                                                                                           ... ........===============================..:,,,,,,.
#                          ,,                                                             ..  .....=.==================================.,,=,,,..
#  '7MMF'     A     '7MF'  db                      .g8"""bgd                      mm     ...........=:~==~~=========================~=.:,...=...
#    'MA     ,MA     ,V                          .dP'     'M                      MM     .......=..=.,:=====~~=======================:=,,.......
#     VM:   ,VVM:   ,V   '7MM  '7Mb,od8  .gP"Ya  dM'       '   ,pW"Wq.  ,pP"Ybd mmMMmm   ......,..:.==.~~~==~~~======================,..=...,...
#      MM.  M' MM.  M'     MM    MM' "' ,M'   Yb MM           6W'   'Wb 8I   '"   MM    .....,..=..=...~===~=~~~~====================.~...=.....
#      'MM A'  'MM A'      MM    MM     8M"""""" MM.    '7MMF'8M     M8 'YMMMa.   MM    ..........:.,=.~=~~=~~~~~~==================~..,=.......
#       :MM;    :MM;       MM    MM     YM.    , 'Mb.     MM  YA.   ,A9 L.   I8   MM    ........=..~...~~~~=~~~~~~~~=================,:+........
#        VF      VF      .JMML..JMML.    'Mbmmd'   '"bmmmdPY   'Ybmd9'  M9mmmP'   'Mbmo......=..~..~..~~~~~~~~~~~~~~~~~=============~~.?,.,.....
                                                                                     ........:.....=.~~~~~~~~~~~~~~~~~~~~~===============.,.....
                                                                                    . .............~~.....~~~~~~~~~~~~~~~~~~=====:.,,.~==,......
                        A Golang Exploitation Framework                              .............,.,.........:+=~~~~~~~~~~~=......,,,,,,=,.....
         __________________________________________________________________         ...............................~~.~~=~........,,,,,,,,+,,...
                                                                                     ..........,...................=.~~~:.........,,,,,,,,,,,...
                                                                                    ................................:.~~~.........,,,,,,,,,~=...
                                                                                       .....::.....................~.=~~~~~.......,,,,,,,,,,.=~.
    " If you think technololy can solve your security problems, then you              .........=..+..........~...:,.~..,~.:~=.+.,,,,,,,,~+:.==..
      don't understand the problems and you don't understand the technology. "         ......=..?=..~=,=,.~..,.~=.......~~,.~~~~~=~~=~=~====...=
                                                                                          ..........,......+..=.+,.......:==.,=====~==,=.,,.....
                                                - Bruce Schneier -                        ...................=,,+:......,.:.=+,++.=,,,.,.,,.=...
                                                                                          .................=..,............=..~.,.~,,,..,,,,,...
                                                                                          ............,......+,.=....,,...+,=~.:~,,,,..,,,......
        %s...........?..~..:.....=,..+.....=..~...,,,,...,......
        %s.............,......,.=.:==.==.=+.~~,..,,,,,,=,.......
                                                                                          ........................=..~..==.,=..:................
        %s......................:...+.~=.,~.....................
        %s........................=.,=,.~:.=...~................
                                                                                          .........................==.,=..=.....:...............
                                                                                          .....................~..=,.=......:...................
`, user, server, cliVer, serVer)

	// Print banners
	fmt.Printf(bannerServerConnection)
}
