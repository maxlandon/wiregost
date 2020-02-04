package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/evilsocket/islazy/tui"

	"github.com/maxlandon/wiregost/client/assets"
	"github.com/maxlandon/wiregost/client/console"
	"github.com/maxlandon/wiregost/client/version"
)

const (
	logFileName = "sliver-client.log"
)

func main() {
	// Print welcome picture
	fmt.Println(tui.Dim("----------------------------------------------------------------" +
		"--------------------------------------------------------------------------------"))
	fmt.Printf(tui.Dim(welcomeToGhost))
	fmt.Println(tui.Dim("----------------------------------------------------------------" +
		"--------------------------------------------------------------------------------"))

	displayVersion := flag.Bool("version", false, "print version number")
	config := flag.String("config", "", "config file")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("v%s\n", version.ClientVersion)
		os.Exit(0)
	}

	// Check configs
	if *config != "" {
		conf, err := assets.ReadConfig(*config)
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(3)
		}
		assets.SaveConfig(conf)
	}

	// Initialize console logging
	appDir := assets.GetRootAppDir()
	logFile := initLogging(appDir)
	defer logFile.Close()

	// Launch session
	console.Start()
}

// Initialize logging
func initLogging(appDir string) *os.File {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile, err := os.OpenFile(path.Join(appDir, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %s", err))
	}
	log.SetOutput(logFile)
	return logFile
}

var welcomeToGhost = `                                                                                                            ,,.,=++============+,               
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
                                                                                          ...........?..~..:.....=,..+.....=..~...,,,,...,......
                                                                                          .............,......,.=.:==.==.=+.~~,..,,,,,,=,.......
                                                                                          ........................=..~..==.,=..:................
                                                                                          ......................:...+.~=.,~.....................
                                                                                          ........................=.,=,.~:.=...~................
                                                                                          .........................==.,=..=.....:...............
                                                                                          .....................~..=,.=......:...................
                                                                                          ......................,~....:...~~....................
                                                                                          ......................................................
`
