package main

// This Go file is an installer for the WireGost client shell, Ghost.

// The workflow is the following:

//  1. Creating a personal user directory.
//	2. Creating base client configuration file.
//  3. Creating all subdirectories and files specific to the client and its user.
//	4. Creating local authentication elements for the client.
//  5. Creating default server paramaters, for first connection.

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/maxlandon/wiregost/internal/session"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
)

var welcomeToGhost = `                                                                                                                 ,,.,=++=======+,                                  
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
	SetUserCreds()
	// Setting up default server
	SetUpDefaultServer()

	// Report and quit
	fmt.Println()
	fmt.Println(tui.Bold("Personal directories, files and authentication information have been setup."))
	fmt.Println(tui.Bold("You can now use the WireGost client shell, Ghost."))
	fmt.Println()

}

// Personal directories and files setup
func PersonalDirFilesSetup() {

	fmt.Println(tui.Yellow(tui.Bold("------------------------------------------------------------------------------------------------------------------------------------------------")))
	fmt.Printf(welcomeToGhost)
	fmt.Println(tui.Yellow(tui.Bold("------------------------------------------------------------------------------------------------------------------------------------------------")))
	fmt.Println()
	fmt.Println(tui.Yellow(tui.Bold("         *********** WireGost Client Setup *********** ")))
	fmt.Println()
	fmt.Println("WireGost is creating personal files and subdirectories in ~/.wiregost/client. They will contain : ")
	fmt.Println()

	// Instantiate a new Config object
	conf := session.NewConfig()

	// Directories
	fmt.Println(tui.Blue(" - Directories -"))
	fmt.Println()
	userdir, _ := fs.Expand(conf.UserDir)
	if fs.Exists(userdir) == false {
		os.MkdirAll(userdir, 0755)
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

	serverCertDir, _ := fs.Expand("~/.wiregost/client/certificates")
	if fs.Exists(serverCertDir) == false {
		os.MkdirAll(serverCertDir, 0755)
		fmt.Println(tui.Blue("  *") + " Servers certificates" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Servers certificates" + tui.Dim(tui.Green(" (Existing)")))
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
func SetUserCreds() (err error) {

	// Instantiate a new User object
	user := session.NewUser()

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
				user.Name = username
				user.PasswordHash = sha256.Sum256(pass)
				user.PasswordHashString = base64.URLEncoding.EncodeToString(user.PasswordHash[:])
				// Save to config file
				credsFile, _ := fs.Expand(user.CredsFile)
				file, _ := os.Create(credsFile)
				defer file.Close()
				defaults, _ := json.MarshalIndent(user, "", "	")
				_, err := file.Write([]byte(defaults))
				if err != nil {
					fmt.Println(err.Error())
				}

				// Report status
				fmt.Println()
				fmt.Println(tui.Blue("Files:"))
				fmt.Printf(tui.Bold("User auth parameters saved to ~/.wiregost/client/.auth file\n"))
				fmt.Printf(tui.Yellow("User name: ")+"%s\n", user.Name)
				fmt.Printf(tui.Yellow("User password hash: ")+"%s\n", user.PasswordHashString)
				fmt.Println()

				break
			}
		}
	}

	return nil
}

func SetUpDefaultServer() error {

	fmt.Println(tui.Dim("------------------------------------------------------------------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Server Configuration -")))
	fmt.Println()
	fmt.Println(" Creating server.conf file in ~/.wiregost/client/")
	fmt.Println()

	// Server Configuration File
	serverFile, _ := fs.Expand("~/.wiregost/client/server.conf")
	if fs.Exists(serverFile) {
		fmt.Println(tui.Blue("  *") + " Configuration file" + tui.Dim(tui.Red(" (Overwritten)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Configuration file" + tui.Dim(tui.Green(" (Created)")))
	}
	servConf, _ := os.Create(serverFile)
	defer servConf.Close()

	// Populate a list of servers with one default server
	var serverList []session.Server
	serverList = []session.Server{
		session.Server{
			IPAddress:   "localhost",
			Port:        7777,
			Certificate: "",
			UserToken:   "",
			FQDN:        "",
			IsDefault:   true,
		}}

	// Status
	fmt.Println()
	fmt.Println("Creating default server parameters (IP, Port), and writing to ~/.wiregost/client/server.conf.")
	// Marshal to JSON
	var jsonData []byte
	jsonData, err := json.MarshalIndent(serverList, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to server configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(serverFile, jsonData, 0755)
		fmt.Println(tui.Green("Server configuration file written."))
	}

	fmt.Println()
	fmt.Println(tui.Bold("IMPORTANT: The default server is not yet ready to connect with TLS support !"))
	fmt.Println()
	fmt.Println(tui.Blue("  1. ") + "Copy the Certificate with extension '.crt' from ~/.wiregost/server/certificates/ to ~/.wiregost/client/certificates/ ")
	fmt.Println(tui.Blue("  2. ") + "Open the ~/.wiregost/client/server.conf file and fill up:")
	fmt.Println(tui.Blue("          * ") + "Fill the Certificate field with ~/.wiregost/client/certificates/your_cert.crt ")
	fmt.Println(tui.Blue("          * ") + "The FQDN field with the " + tui.Yellow("Common Name ") + "or " + tui.Yellow("server FQDN"))
	fmt.Println("            that you filled when running the server setup file. This name can be empty,")
	fmt.Println("            as long as they are identical", " in the Certificate and the server.conf file).")
	fmt.Println()
	fmt.Println(tui.Dim("------------------------------------------------------------------------------------------------------------------------------------------------"))

	return nil
}
