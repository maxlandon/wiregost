package main

// This Go file is an installer for the WireGost client shell, Ghost.

// The workflow is the following:

//  1. Creating a personal user directory.
//	2. Creating base client configuration file.
//  3. Creating all subdirectories and files specific to the client and its user.
//	4. Creating local authentication elements for the client.

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/maxlandon/wiregost/internal/session"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
)

var welcomeWireGost = `                                                                                                                 ,,.,=++=======+,                                  
                                                                                                         ..~====================+..                            
                                                                                                       ,+===========================..                         
                                                                                                      :===============================                         
#                                                                                           .  .......+===============================:.                       
                                                                                           ... ........===============================..:,,,,,,. 
#                          ,,                                                             ..  .....=.==================================.,,=,,,..               
#  '7MMF'     A     '7MF'  db                      .g8"""bgd                      mm    . . . ......+..==============================~,,=,,,,...               
#    'MA     ,MA     ,V                          .dP'     'M                      MM     . ........:.================================~+,,,~,,,,.               
#     VM:   ,VVM:   ,V   '7MM  '7Mb,od8  .gP"Ya  dM'       '   ,pW|Wq.  ,pP"Ybd mmMMmm   . .. ...,.:~,=~~============================+,~=.,.,...               
#      MM.  M' MM.  M'     MM    MM' "' ,M'   Yb MM           6W'   'Wb 8I   '"   MM     ..........=.:~=~~===========================+~:,.:,,...               
#      'MM A'  'MM A'      MM    MM     8M"""""" MM.    '7MMF'8M     M8 'YMMMa.   MM     ...........=:~==~~=========================~=.:,...=...               
#       :MM;    :MM;       MM    MM     YM.    , 'Mb.     MM  YA.   ,A9 L.   I8   MM     .......=..=.,:=====~~=======================:=,,.......               
#        VF      VF      .JMML..JMML.    'Mbmmd'   '"bmmmdPY   'Ybmd9'  M9mmmP'   'Mbmo  ......,..:.==.~~~==~~~======================,..=...,...               
#                                                                                       .....,..=..=...~===~=~~~~====================.~...=.....               
                                                                                        ..........:.,=.~=~~=~~~~~~==================~..,=.......               
                        A Golang Exploitation Framework                                 ........=..~...~~~~=~~~~~~~~=================,:+........               
       __________________________________________________________________               .......:..=..=.~=~==~~~~~~~~~================...,~......               
                                                                                       ......=..~..~..~~~~~~~~~~~~~~~~~=============~~.?,.,.....               
                                                                                      .........=.....==~~~~~~~~~~~~~~~~~==============+,,+......               
                                                                                     ........:.....=.~~~~~~~~~~~~~~~~~~~~~===============.,.....               
    " If you think technololy can solve your security problems, then you            . .............~~.....~~~~~~~~~~~~~~~~~~=====:.,,.~==,......               
      don't understand the problems and you don't understand the technology. "       .............,.,.........:+=~~~~~~~~~~~=......,,,,,,=,.....               
                                                                                    ...............................~~.~~=~........,,,,,,,,+,,...               
                                                - Bruce Schneier -                   ..........,...................=.~~~:.........,,,,,,,,,,,...               
                                                                                    ................................:.~~~.........,,,,,,,,,~=...               
                                                                                       .....::.....................~.=~~~~~.......,,,,,,,,,,.=~.               
                                                                                      .........=..+..........~...:,.~..,~.:~=.+.,,,,,,,,~+:.==..               
                                                                                       ......=..?=..~=,=,.~..,.~=.......~~,.~~~~~=~~=~=~====...=               
                                                                                          ..........,......+..=.+,.......:==.,=====~==,=.,,.....               
                                                                                          ...................=,,+:......,.:.=+,++.=,,,.,.,,.=...               
                                                                                          .................=..,............=..~.,.~,,,..,,,,,...               
                                                                                          ............,......+,.=....,,...+,=~.:~,,,,..,,,......               
                                                                                          ...........?..~..:.....=,..+.....=..~...,,,,...,......               
                                                                                          .............,......,.=.:==.==.=+.~~,..,,,,,,=,.......               
                                                                                          .......................==.++,=+,,~...:,.,,,.,.........               
                                                                                          .....................+=.~=:,~~.~=.==...,,,,,..........               
                                                                                          .....................=,~+,=+,~=..=,..:..,.............               
                                                                                          ...................:.,~.,~~,==.++.=:...?.....,........               
                                                                                          ..................~....=:..=..=..==..=................               
                                                                                          .............~..+...,,+..:+.==.==,.~~,.=..:...........               
                                                                                          ........................=..~..==.,=..:................               
                                                                                          ......................:...+.~=.,~.....................               
                                                                                          ........................=.,=,.~:.=...~................               
                                                                                          .........................==.,=..=.....:...............               
                                                                                          .....................~..=,.=......:...................               
                                                                                          ......................,~....:...~~....................               
                                                                                          ......................................................               
`

func main() {

	// Setting up personal directories and files
	PersonalDirFilesSetup()
	// Setting aup client-side local authentication
	SetClientAuth()

	// Report and quit
	fmt.Println()
	fmt.Println(tui.Bold("Personal directories, files and authentication information have been setup. You can now use the WireGost client shell, Ghost."))

}

// Personal directories and files setup
func PersonalDirFilesSetup() {

	fmt.Println(tui.Yellow(tui.Bold("------------------------------------------------------------------------------------------------------------------------------------------------")))
	fmt.Printf(welcomeWireGost)
	fmt.Println(tui.Yellow(tui.Bold("------------------------------------------------------------------------------------------------------------------------------------------------")))
	fmt.Println()
	fmt.Println(tui.Yellow(tui.Bold("         *********** WireGost Client Setup *********** ")))
	fmt.Println()
	fmt.Println("WireGost is creating personal files and subdirectories in ~/.wiregost. They will contain : ")
	fmt.Println()

	// Instantiate a new Config object
	conf := session.NewConfig()

	// Directories
	fmt.Println(tui.Blue(" - Directories -"))
	userdir, _ := fs.Expand(conf.UserDir)
	if fs.Exists(userdir) == false {
		os.Mkdir(userdir, 0755)
		fmt.Println(tui.Blue("  *") + " User directory" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " User directory" + tui.Dim(tui.Green(" (Existing)")))
	}

	logdir, _ := fs.Expand(conf.LogDir)
	if fs.Exists(logdir) == false {
		os.Mkdir(logdir, 0755)
		fmt.Println(tui.Blue("  *") + " Logs" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Logs" + tui.Dim(tui.Green(" (Existing)")))
	}

	moduledir, _ := fs.Expand(conf.ModulesDir)
	if fs.Exists(moduledir) == false {
		os.Mkdir(moduledir, 0755)
		fmt.Println(tui.Blue("  *") + " Modules" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Modules" + tui.Dim(tui.Green(" (Existing)")))
	}

	payloaddir, _ := fs.Expand(conf.PayloadsDir)
	if fs.Exists(payloaddir) == false {
		os.Mkdir(payloaddir, 0755)
		fmt.Println(tui.Blue("  *") + " Payloads" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Payloads" + tui.Dim(tui.Green(" (Existing)")))
	}

	resourcedir, _ := fs.Expand(conf.ResourceDir)
	if fs.Exists(resourcedir) == false {
		os.Mkdir(resourcedir, 0755)
		fmt.Println(tui.Blue("  *") + " Resources" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Resources" + tui.Dim(tui.Green(" (Existing)")))
	}

	workspacedir, _ := fs.Expand(conf.WorkspaceDir)
	if fs.Exists(workspacedir) == false {
		os.Mkdir(workspacedir, 0755)
		fmt.Println(tui.Blue("  *") + " Workspaces" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Workspaces" + tui.Dim(tui.Green(" (Existing)")))
	}

	exportdir, _ := fs.Expand(conf.ExportDir)
	if fs.Exists(exportdir) == false {
		os.Mkdir(exportdir, 0755)
		fmt.Println(tui.Blue("  *") + " Exports" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Exports" + tui.Dim(tui.Green(" (Existing)")))
	}
	fmt.Println()

	// Files
	fmt.Println(tui.Blue(" - Files - "))
	configfile, _ := fs.Expand(conf.UserConfigFile)
	if fs.Exists(configfile) {
		fmt.Println(tui.Blue("  *") + " Configuration file" + tui.Dim(tui.Red(" (Overwritten)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Configuration file" + tui.Dim(tui.Green(" (Created)")))
	}
	file, _ := os.Create(configfile)
	defer file.Close()
	defaults, _ := json.MarshalIndent(conf, "", "	")
	file.Write(defaults)

	historyfile, _ := fs.Expand(conf.HistoryFile)
	if fs.Exists(historyfile) {
		fmt.Println(tui.Blue("  *") + " Command History file " + tui.Dim(tui.Red(" (Overwritten)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Command History file " + tui.Dim(tui.Green(" (Created)")))
	}
	histfile, _ := os.Create(historyfile)
	defer histfile.Close()
}

// Setup client-side local authentication
func SetClientAuth() (err error) {

	// Instantiate a new Auth object
	auth := session.NewAuth()

	fmt.Println(tui.Dim("------------------------------------------------------------------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Authentication -")))
	fmt.Println()
	fmt.Println("Please follow prompt for creating ID and authentication parameters.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	username := ""
	password := ""

	// Read user name input
	fmt.Printf(tui.Bold("Enter user name: "))
	for {
		input, _ := reader.ReadString('\n')
		username = strings.TrimRight(input, "\n")
		if username != "" {
			fmt.Printf("User Name: %s \n", username)
			break
		}
	}

	// Read user password input
	for {
		fmt.Printf(tui.Bold("Enter user password: "))
		pass, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		password = string(pass)
		if err == nil {
			fmt.Printf(tui.Bold("Confirm user password: "))
			fmt.Println()
			con, _ := terminal.ReadPassword(int(syscall.Stdin))
			confirm := string(con)
			if password != confirm {
				fmt.Println(tui.Red("Passwords mismatch. Retry"))
				fmt.Println()
			} else {
				auth.UserName = username
				auth.PasswordHash = sha256.Sum256(pass)
				auth.PasswordHashString = base64.URLEncoding.EncodeToString(auth.PasswordHash[:])
				// Save to config file
				authFile, _ := fs.Expand(auth.UserAuthFile)
				file, _ := os.Create(authFile)
				defer file.Close()
				defaults, _ := json.MarshalIndent(auth, "", "	")
				_, err := file.Write([]byte(defaults))
				if err != nil {
					fmt.Println(err.Error())
				}

				// Report status
				fmt.Println()
				fmt.Println(tui.Blue("Files:"))
				fmt.Printf(tui.Bold("User auth parameters saved to ~/.wiregost/.auth file\n"))
				fmt.Printf(tui.Yellow("User name: ")+"%s\n", username)
				fmt.Printf(tui.Yellow("User password hash: ")+"%s\n", auth.PasswordHashString)
				fmt.Println(tui.Dim("------------------------------------------------------------------"))

				break
			}
		}
	}

	return nil
}
