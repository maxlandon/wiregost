package session

// This file contains all the code used to manage servers and their connections.
// The following are included:

// LOCAL -----------
//	- Types for Servers, list of saved Servers, currently connected User
//	- Paths to user configuration files
//  - Functions for managing the list of saved Servers, getting current
//	  and default ones

// RPC -------------
//	- Types for RPC security
//	- Functions for RPC security
//	- RPC functions for issuing commands remotely to the connected server.
//  - RPC functions reserved to administrators

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/maxlandon/wiregost/internal/server/core"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	IPAddress   string
	Port        int
	Certificate string
	UserToken   string
	FQDN        string
	IsDefault   bool
}

type ServerManager struct {
	SavedServers  []Server
	CurrentServer Server
	creds         credentials.TransportCredentials
	auth          Authentication
	connected     bool
}

var serverFile = "~/.wiregost/client/server.conf"

type Authentication struct {
	Login    string
	Password string
}

//----------------------------------------------------------
// LOCAL FUNCTIONS

func NewServerManager(user *User) *ServerManager {
	sv := &ServerManager{connected: false}

	// Load saved servers
	sv.GetServerList()
	sv.GetDefaultServer()

	// Set up credentials
	sv.creds, err = credentials.NewClientTLSFromFile(sv.CurrentServer.Certificate, sv.CurrentServer.FQDN)
	if err != nil {
		fmt.Println(tui.Red("Could not load TLS certificates."))
	}
	// Setup auth
	sv.auth = Authentication{
		Login:    user.Name,
		Password: user.PasswordHashString,
	}
	// Connect to default server (CHANGE THIS WHEN CONNECT FUNCTION IS DONE)
	sv.ConnectToServer(user, sv.CurrentServer)

	return sv
}

func (sv *ServerManager) GetServerList() error {
	serverList := []Server{}
	path, _ := fs.Expand(serverFile)
	if !fs.Exists(path) {
		fmt.Println(tui.Red("Configuration file not found: check for issues," +
			" or run the configuration script again"))
		os.Exit(1)
	} else {
		configBlob, _ := ioutil.ReadFile(path)
		json.Unmarshal(configBlob, &serverList)
	}

	// Format certificate path for each server, add server to ServerManager
	for _, i := range serverList {
		i.Certificate, _ = fs.Expand(i.Certificate)
		sv.SavedServers = append(sv.SavedServers,
			Server{IPAddress: i.IPAddress,
				Port:        i.Port,
				Certificate: i.Certificate,
				FQDN:        i.FQDN,
				IsDefault:   i.IsDefault})
	}
	return nil
}

func (sv *ServerManager) WriteServerList() error {

	return nil
}

func (sv *ServerManager) AddServer() error {

	return nil
}

func (sv *ServerManager) DeleteServer() error {

	return nil
}

func (sv *ServerManager) SetDefaultServer() error {

	return nil
}

func (sv *ServerManager) GetDefaultServer() error {
	for _, i := range sv.SavedServers {
		if i.IsDefault == true {
			sv.CurrentServer = i
			break
		}
	}
	return nil
}

//----------------------------------------------------------
// RPC SECURITY

// Authentication
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

//----------------------------------------------------------
// RPC SERVICES

func (sv *ServerManager) DialRPC() (*grpc.ClientConn, error) {
	// Initiate connection with the server
	conn, err := grpc.Dial(sv.CurrentServer.IPAddress+":"+strconv.Itoa(sv.CurrentServer.Port),
		grpc.WithTransportCredentials(sv.creds),
		grpc.WithPerRPCCredentials(&sv.auth))
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	return conn, err
}

func (sv *ServerManager) ConnectToServer(user *User, server Server) error {
	// If already connected to a server, disconnect
	if sv.connected == true {
		sv.DisconnectFromServer()
	}
	// Change CurrentServer
	sv.CurrentServer = server

	conn, err := sv.DialRPC()
	defer conn.Close()
	client := core.NewUserManagerClient(conn)

	request := &core.ConnectRequest{}
	response, err := client.ConnectUser(context.Background(), request)
	if err != nil {
		log.Println(tui.Red("Error: could not get response from server."))
		log.Println(tui.Red(err.Error()))
		return nil
	}

	if response.Clearance == "clear" && response.Admin == true {
		sv.connected = true
		log.Printf("Connected as " + tui.Bold(tui.Yellow(user.Name)+" (Administrator rights)"))
		log.Printf("Server at " + sv.CurrentServer.IPAddress + ":" + strconv.Itoa(sv.CurrentServer.Port) + " (FQDN: " +
			sv.CurrentServer.FQDN + ", default: " + strconv.FormatBool(sv.CurrentServer.IsDefault) + ")")
	}
	if response.Clearance == "clear" && response.Admin == false {
		sv.connected = true
		log.Printf("Connected as %s", tui.Bold(tui.Yellow(user.Name)))
		log.Printf("Server at "+sv.CurrentServer.IPAddress, ":"+strconv.Itoa(sv.CurrentServer.Port)+
			"(FQDN: "+sv.CurrentServer.FQDN+", default: "+strconv.FormatBool(sv.CurrentServer.IsDefault)+")")
	}
	if response.Clearance == "reg" && response.Admin == false {
		sv.connected = true
		log.Printf(tui.Green("First connection of user ")+tui.Bold(tui.Yellow(user.Name)),
			" : User and password are now registered in the server database.")
		fmt.Println()
		log.Printf("Connected as %s", tui.Bold(tui.Yellow(user.Name)))
		log.Printf("Server at "+sv.CurrentServer.IPAddress, ":"+strconv.Itoa(sv.CurrentServer.Port)+
			"(FQDN: "+sv.CurrentServer.FQDN+", default: "+strconv.FormatBool(sv.CurrentServer.IsDefault)+")")
	}
	if response.Clearance == "reg" && response.Admin == true {
		sv.connected = true
		log.Printf(tui.Green("First connection of user ")+tui.Bold(tui.Yellow(user.Name)),
			" : User and password are now registered in the server database.")
		fmt.Println()
		log.Printf("Connected as " + tui.Bold(tui.Yellow(user.Name)+" (Administrator rights)"))
		log.Printf("Server at "+sv.CurrentServer.IPAddress, ":"+strconv.Itoa(sv.CurrentServer.Port)+
			"(FQDN: "+sv.CurrentServer.FQDN+", default: "+strconv.FormatBool(sv.CurrentServer.IsDefault)+")")
	}
	if response.Clearance == "none" {
		log.Printf(tui.Red("No user ") + tui.Bold(tui.Yellow(user.Name)) + " in database. Connection aborted")
	}

	return nil
}

func (sv *ServerManager) DisconnectFromServer() error {

	conn, err := sv.DialRPC()
	defer conn.Close()
	client := core.NewUserManagerClient(conn)

	request := &core.DisconnectRequest{}
	_, err = client.DisconnectUser(context.Background(), request)
	if err != nil {
		log.Println(tui.Red("Error: could not get response from server."))
		log.Println(tui.Red(err.Error()))
		return nil
	}
	return nil
}
