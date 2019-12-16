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
}

var serverFile = "~/.wiregost/client/server.conf"

type Authentication struct {
	Login    string
	Password string
}

//----------------------------------------------------------
// LOCAL FUNCTIONS

func NewServerManager(user *User) *ServerManager {
	sv := &ServerManager{}

	// Load saved servers
	sv.GetServerList()
	sv.GetDefaultServer()

	// Set up credentials
	sv.creds, err = credentials.NewClientTLSFromFile(sv.CurrentServer.Certificate, sv.CurrentServer.FQDN)
	if err != nil {
		fmt.Println(tui.Red("Could not load TLS certificates."))
		fmt.Println("here")
	}

	// Setup auth
	sv.auth = Authentication{
		Login:    user.Name,
		Password: user.PasswordHashString,
	}

	// Connect to default server (CHANGE THIS WHEN CONNECT FUNCTION IS DONE)
	sv.RegisterUserToServer(user)

	return sv
}

func (sv *ServerManager) GetServerList() error {
	serverList := []Server{}
	path, _ := fs.Expand(serverFile)
	if !fs.Exists(path) {
		fmt.Println(tui.Red("Configuration file not found: check for issues, or run the configuration script again"))
		os.Exit(1)
	} else {
		configBlob, _ := ioutil.ReadFile(path)
		json.Unmarshal(configBlob, &serverList)
		fmt.Println(tui.Dim("Configuration file loaded."))
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
func (sv *ServerManager) RegisterUserToServer(user *User) {
	var conn *grpc.ClientConn

	// Initiate connection with the server
	conn, err = grpc.Dial(sv.CurrentServer.IPAddress+":"+strconv.Itoa(sv.CurrentServer.Port),
		grpc.WithTransportCredentials(sv.creds))
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := core.NewUserManagerClient(conn)

	request := &core.RegisterRequest{Name: user.Name, Hash: user.PasswordHashString}
	response, err := client.RegisterUser(context.Background(), request)
	if err != nil {
		log.Println(err)
		return
	}
	if response.Registered == true && response.Error == "" {
		log.Printf(tui.Green("Client is now registered. PasswordHash is saved to database."))
		log.Printf("You can now connect to the DB, the next time you use the Ghost client.")
	}
	if response.Registered == true && response.Error != "" {
		log.Printf(response.Error)
		if response.Registered == false && response.Error != "" {
			log.Printf(response.Error)
		}
		if response.Registered == false && response.Error == "" {
			log.Println(tui.Red("Error: user is either not registered, or the server has mishandled the request/database"))
		}
	}
}

func (sv *ServerManager) ExampleFuncWithFullSecurity(user *User) {
	var conn *grpc.ClientConn

	// Create TLS Credentials
	creds, err := credentials.NewClientTLSFromFile(sv.CurrentServer.Certificate, sv.CurrentServer.FQDN)
	if err != nil {
		fmt.Println(tui.Red("Could not load TLS certificates."))
	}

	// Set login/pass
	auth := Authentication{
		Login:    user.Name,
		Password: user.PasswordHashString,
	}

	// Initiate connection with the server
	conn, err = grpc.Dial(sv.CurrentServer.IPAddress+":"+strconv.Itoa(sv.CurrentServer.Port),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := core.NewUserManagerClient(conn)

	request := &core.RegisterRequest{Name: user.Name, Hash: user.PasswordHashString}
	response, err := client.RegisterUser(context.Background(), request)
	if err != nil {
		log.Println(tui.Red("Error: user is either not registered, or the server has mishandled the request/database"))
		return
	}
	if response.Registered == true {
		log.Printf(tui.Green("Client is now registered !"))
	}
}

func (sv *ServerManager) ConnectToServer(user *User, server Server) error {

	return nil
}
