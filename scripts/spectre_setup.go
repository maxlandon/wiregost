package main

// This Go file is an installer for the WireGost server, Spectre.

// The workflow is the following:

// 1. Create database.
// 2. Create database tables.
// 3. Create user name, with admin rights given

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/maxlandon/wiregost/internal/server/core"
	"github.com/maxlandon/wiregost/internal/server/db"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var welcomeToSpectre = `                                                                                                            ,,.,=++============+,       
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

func main() {

	fmt.Println(tui.Dim(tui.Yellow("----------------------------------------------------------" +
		"--------------------------------------------------------------------------------------")))
	fmt.Printf(welcomeToSpectre)
	fmt.Println(tui.Dim(tui.Yellow("----------------------------------------------------------" +
		"--------------------------------------------------------------------------------------")))
	fmt.Println()
	fmt.Println(tui.Yellow(tui.Bold("         *********** WireGost Server Setup *********** ")))
	fmt.Println()

	// Creating Directories
	CreateDirectories()

	// Database Setup
	DatabaseSetup()

	// Creating tables
	CreateUsersTable()

	// CreateDefaultUser
	CreateDefaultUser()

	// Create Certificates
	CreateCertificates()

}

// ----------------------------------------------------
// Directories Setup
func CreateDirectories() {

	// Directories
	fmt.Println(tui.Bold(tui.Blue(" - Directories -")))
	serverDir, _ := fs.Expand("~/.wiregost/server")
	if fs.Exists(serverDir) == false {
		os.MkdirAll(serverDir, 0755)
		fmt.Println()
		fmt.Println(tui.Blue("  *") + " Server directory" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println()
		fmt.Println(tui.Blue("  *") + " Server directory" + tui.Dim(tui.Green(" (Existing)")))
	}
	certsDir, _ := fs.Expand("~/.wiregost/server/certificates")
	if fs.Exists(certsDir) == false {
		os.MkdirAll(certsDir, 0755)
		fmt.Println(tui.Blue("  *") + " Certificates directory" + tui.Dim(tui.Green(" (Created)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Certificates directory" + tui.Dim(tui.Green(" (Existing)")))
	}
	fmt.Println()
}

// ----------------------------------------------------
// Database Setup

func DatabaseSetup() {

	fmt.Println(tui.Dim("--------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Database -")))
	fmt.Println()

	// Input Dabase name, user and password
	reader := bufio.NewReader(os.Stdin)
	dbUser := ""
	dbName := ""
	dbPass := ""

	fmt.Printf(tui.Bold("Enter database user name: "))
	input, _ := reader.ReadString('\n')
	dbUser = strings.TrimRight(input, "\n")
	if dbUser != "" {
		fmt.Printf("Db user name: %s \n", dbUser)
		fmt.Println()
	}
	if dbUser == "" {
		fmt.Println(tui.Bold("Db User cannot be empty, please provide one:"))
		input, _ = reader.ReadString('\n')
		dbUser = strings.TrimRight(input, "\n")
		fmt.Printf("Db user name: %s \n", dbUser)
		fmt.Println()
	}

	fmt.Printf(tui.Bold("Enter database name: "))
	input, _ = reader.ReadString('\n')
	dbName = strings.TrimRight(input, "\n")
	if dbName != "" {
		fmt.Printf("Db name: %s \n", dbName)
		fmt.Println()
	}
	if dbName == "" {
		fmt.Println(tui.Bold("Db Name cannot be empty, please provide one:"))
		input, _ = reader.ReadString('\n')
		dbName = strings.TrimRight(input, "\n")
		fmt.Printf("Db name: %s \n", dbName)
		fmt.Println()
	}

	fmt.Printf(tui.Bold("Enter database password: "))
	input, _ = reader.ReadString('\n')
	dbPass = strings.TrimRight(input, "\n")
	if dbPass != "" {
		fmt.Printf("Db password: %s \n", dbPass)
		fmt.Println()
	}
	if dbPass == "" {
		fmt.Println(tui.Bold("Password cannot be empty, please provide one:"))
		input, _ = reader.ReadString('\n')
		dbPass = strings.TrimRight(input, "\n")
		fmt.Printf("Db password: %s \n", dbPass)
		fmt.Println()
	}

	userPhrase := fmt.Sprintf("CREATE USER %s WiTH PASSWORD '%s';", dbUser, dbPass)
	dbPhrase := fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s;", dbName, dbUser)

	cmd := exec.Command("psql", "-U", "postgres",
		"-c",
		userPhrase,
		"-c",
		dbPhrase)

	// Error handling
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(tui.Red(fmt.Sprint(err) + ": " + stderr.String()))
	}

	// Create database configuration file and populate
	fmt.Println("WireGost default database created :")
	fmt.Printf("User: %s \n", tui.Yellow(dbUser))
	fmt.Printf("Password: %s \n", tui.Yellow(dbPass))
	fmt.Printf("Database: %s \n", tui.Yellow(dbName))
	fmt.Println()

	dbFile, _ := fs.Expand("~/.wiregost/server/database.conf")
	if fs.Exists(dbFile) {
		fmt.Println(tui.Blue("  *") + " Database Configuration file" + tui.Dim(tui.Red(" (Overwritten)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Database Configuration file" + tui.Dim(tui.Green(" (Created)")))
	}
	dbConf, _ := os.Create(dbFile)
	defer dbConf.Close()

	type DbCreds struct {
		User     string
		Database string
		Password string
	}
	opts := DbCreds{
		User:     dbUser,
		Database: dbName,
		Password: dbPass,
	}

	var jsonData []byte
	jsonData, err = json.MarshalIndent(opts, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to database configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(dbFile, jsonData, 0755)
		fmt.Println("Populated dabatase configuration file with credentials")
	}
}

// This function will need to be renamed "CreateDefaultTables",
// and the code for creating all other tables be added to it,
// so that all tables are created in the same function.
func CreateUsersTable() {

	// Presentation
	fmt.Println()
	fmt.Println("Creating database tables :")
	fmt.Println()

	opts := &pg.Options{
		User:     "wiregost",
		Password: "wiregost",
		Database: "wiregost_db",
	}
	var database *pg.DB = pg.Connect(opts)
	_, health := database.Exec("SELECT 1")

	// Quit if failed to connect
	if health != nil {
		log.Println(tui.Bold(tui.Red("Failed to connect to database")))
		log.Println(tui.Bold(tui.Red("Please check that Postgresql is running and can be accessed")))
		os.Exit(1)
	} else {
		log.Println(tui.Dim(tui.Green("Connection to database successful.")))
	}

	// Create User Table
	options := &orm.CreateTableOptions{
		//IfNotExists: true,
	}
	createErr := database.CreateTable(&db.User{}, options)
	if createErr != nil {
		log.Printf("Error while creating table Users. Reason: %v\n", createErr)
	} else {
		log.Println(tui.Green("Created table: Users"))
	}

	// Create other tables
	createErr = database.CreateTable(&db.Workspace{}, options)
	if createErr != nil {
		log.Printf("Error while creating table Workspaces. Reason: %v\n", createErr)
	} else {
		log.Println(tui.Green("Created table: Workspaces"))
	}

	// Close DB
	closeErr := database.Close()
	if closeErr != nil {
		log.Printf("Error while closing the connection. Reason: %v\n", closeErr)
		os.Exit(1)
	} else {
		log.Println(tui.Dim(tui.Green("Connection closed successfully.")))
	}

	fmt.Println()
}

// ----------------------------------------------------
// Default User Setup

func CreateDefaultUser() {

	fmt.Println(tui.Dim("--------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Default User -")))
	fmt.Println()
	fmt.Println("Please follow prompt for creating a default User.")
	fmt.Println("This user will be registered in the WireGost server database, with administrator rights.")
	fmt.Println()
	fmt.Println("During WireGost client setup" + tui.Bold(" YOU MUST ENTER THE SAME USER NAME !"))
	fmt.Println("This will allow you to fully register your user, using the client shell (password, hash, etc).")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	username := ""

	// Read user name input
	fmt.Printf(tui.Bold("Enter user name: "))
	for {
		input, _ := reader.ReadString('\n')
		username = strings.TrimRight(input, "\n")
		if username != "" {
			fmt.Printf("User Name: %s \n", username)
			fmt.Println()
			break
		}
	}
	user := &db.User{Name: username, Administrator: true}

	// Save user in database
	opts := &pg.Options{
		User:     "wiregost",
		Password: "wiregost",
		Database: "wiregost_db",
	}
	var database *pg.DB = pg.Connect(opts)

	// Insert
	err := database.Insert(user)
	if err != nil {
		fmt.Println(tui.Red("A problem happened, user was not registered into WireGost database."))
	} else {
		fmt.Println(tui.Green("User ") + tui.Bold(tui.Yellow(username)) +
			tui.Green(" has been registered in the database."))
	}

	// Create default workspace
	user = new(db.User)
	err = database.Model(user).Where("name = ?", "para").Select()

	workspace := &db.Workspace{
		Name:      "default",
		OwnerID:   user.Id,
		CreatedAt: time.Now().String(),
	}

	err = database.Insert(workspace)
	if err != nil {
		fmt.Println(tui.Dim("Could not create default workspace."))
	}
}

// ----------------------------------------------------
// Certificates Creation
func CreateCertificates() error {

	userDir, _ := os.UserHomeDir()

	fmt.Println()
	fmt.Println(tui.Dim("--------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Transport Layer Security (SSL/TLS) -")))
	fmt.Println()

	fmt.Println("Creating SSL Certificates and private key for the server (in ~/.wiregost/server/certificates)")
	fmt.Println()
	fmt.Printf(tui.Bold("Enter a name for the Certificates (without file extension): "))

	reader := bufio.NewReader(os.Stdin)
	certName := ""

	input, _ := reader.ReadString('\n')
	certName = strings.TrimRight(input, "\n")
	if certName != "" {
		fmt.Printf("Certificate name: %s \n", certName)
		fmt.Println()
	}
	if certName == "" {
		fmt.Println(tui.Bold("Name cannot be empty, please provide one:"))
		input, _ = reader.ReadString('\n')
		certName = strings.TrimRight(input, "\n")
		fmt.Printf("Certificate name: %s \n", certName)
		fmt.Println()
	}

	// Create private key
	fmt.Println("Generating Private Key")
	cmd := exec.Command("openssl", "genrsa", "-out", userDir+"/.wiregost/server/certificates/"+certName+".key", "2048")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: failed to generate private key. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		fmt.Println(tui.Green("Generated Private Key 'spectre.key'"))
	}

	// Create Certificate
	fmt.Println()
	fmt.Println(tui.Dim(tui.Blue("Generating Certificate from Private Key. (Expiration date: 3650 days)")))
	cmd = exec.Command("openssl", "req", "-new", "-x509", "-sha256", "-key",
		userDir+"/.wiregost/server/certificates/"+certName+".key", "-out",
		userDir+"/.wiregost/server/certificates/"+certName+".crt", "-days", "3650")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error: failed to generate Certificate. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		fmt.Println(tui.Green("Generated Certificate 'spectre.crt'"))
	}

	// Generate self-signed Certificate
	fmt.Println()
	fmt.Println(tui.Dim(tui.Blue("Generating self-signed Certificate from Private Key")))
	cmd = exec.Command("openssl", "req", "-new", "-sha256", "-key",
		userDir+"/.wiregost/server/certificates/"+certName+".key", "-out",
		userDir+"/.wiregost/server/certificates/"+certName+".csr")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: failed to generate self-signed Certificate. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		fmt.Println(tui.Green("Generated self-signed Certificate 'spectre.csr'"))
	}

	// Sign Private Key with self-signed Certificate
	fmt.Println()
	fmt.Println(tui.Dim(tui.Blue("Signing Private Key with self-signed Certificate. (Expiration date: 3650 days)")))
	cmd = exec.Command("openssl", "x509", "-req", "-sha256", "-in",
		userDir+"/.wiregost/server/certificates/"+certName+".csr", "-signkey",
		userDir+"/.wiregost/server/certificates/"+certName+".key", "-out",
		userDir+"/.wiregost/server/certificates/"+certName+".crt", "-days", "3650")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: failed to sign Private Key with self-signed Certificate. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		fmt.Println(tui.Green("Signed Private Key with self-signed Certificate 'spectre.csr'"))
		fmt.Println()
	}

	// Create Server configuration file, here only because we need the CertName entered previously,
	// For populating default server configuration.
	fmt.Println(tui.Dim("--------------------------------------------------------------------------------------"))
	fmt.Println()
	fmt.Println(tui.Bold(tui.Blue(" - Server Configuration File -")))
	fmt.Println()
	serverFile, _ := fs.Expand("~/.wiregost/server/server.conf")
	if fs.Exists(serverFile) {
		fmt.Println(tui.Blue("  *") + " Server Configuration file" + tui.Dim(tui.Red(" (Overwritten)")))
	} else {
		fmt.Println(tui.Blue("  *") + " Server Configuration file" + tui.Dim(tui.Green(" (Created)")))
	}
	servConf, _ := os.Create(serverFile)
	defer servConf.Close()

	serv := core.ClientRPC{
		Protocol:  "tcp",
		IpAddress: "localhost",
		Port:      7777,
		CertPath:  "~/.wiregost/server/certificates/" + certName + ".crt",
		KeyPath:   "~/.wiregost/server/certificates/" + certName + ".key",
	}

	var jsonData []byte
	jsonData, err = json.MarshalIndent(serv, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to server configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(serverFile, jsonData, 0755)
		fmt.Println()
		fmt.Println(tui.Bold("Default server configuration:"))
		fmt.Println("Protocol: tcp")
		fmt.Println("IP address: localhost")
		fmt.Println("Port: 7777")
		fmt.Println()
		fmt.Println("Certificate & key path: ~/.wiregost/server/certificates/")
		fmt.Println()
	}

	fmt.Println()
	fmt.Println(tui.Bold("The Certificate 'spectre.cert' will also be needed for the Ghost client."))
	fmt.Println(tui.Bold("After finishing the server AND the client setup (and before first connection)"))
	fmt.Println(tui.Bold("please copy the Certificate into the client's personal directory (~/.wiregost/client/server_certs)."))
	fmt.Println()
	fmt.Println(tui.Bold("You will also need the ") + tui.Bold(tui.Yellow("Common Name/ FQDN")) +
		" that you used above when creating your certificates.")

	return nil
}
