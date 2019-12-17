package core

// This file contains the code for the RPC facing the WireGost client.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/maxlandon/wiregost/internal/server/db"

	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var userHomeDir, _ = os.UserHomeDir()
var serverFile = userHomeDir + "/.wiregost/server/server.conf"

type ClientRPC struct {
	Protocol  string
	IpAddress string
	Port      int
	CertPath  string
	KeyPath   string
	creds     credentials.TransportCredentials
	opts      []grpc.ServerOption
	server    *grpc.Server
	Users     []db.User
}

var clearanceCtx string
var adminCtx bool

func NewClientRPC() *ClientRPC {
	serv := &ClientRPC{}

	return serv
}

func (serv *ClientRPC) Start() error {

	// Prepare listener
	lis, err := net.Listen(serv.Protocol, fmt.Sprintf("%s:%d", "localhost", serv.Port))
	if err != nil {
		log.Fatalf("%s Failed to listen on port %d: %v", tui.RED, serv.Port, err)
	}

	// Start the server
	if err := serv.server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}

	return nil
}

func (serv *ClientRPC) LoadConfig() error {

	// Load configuration
	fmt.Println(tui.Dim("Personal directory found."))
	path, _ := fs.Expand(serverFile)
	conf := ClientRPC{}
	// If config file doesn't exist, exit the client
	if !fs.Exists(path) {
		fmt.Println(tui.Red("Configuration file not found: check for issues, or run the server configuration script again"))
		os.Exit(1)
		// If config file is found, parse it.
	} else {
		configBlob, _ := ioutil.ReadFile(path)
		json.Unmarshal(configBlob, &conf)
		fmt.Println(tui.Dim("Configuration file loaded."))
	}

	// Format conf and fill ClientRPC
	serv.Protocol = conf.Protocol
	serv.IpAddress = conf.IpAddress
	serv.Port = conf.Port
	serv.CertPath, _ = fs.Expand(conf.CertPath)
	serv.KeyPath, _ = fs.Expand(conf.KeyPath)

	// Load TLS Credentials
	var err error
	serv.creds, err = credentials.NewServerTLSFromFile(serv.CertPath, serv.KeyPath)
	if err != nil {
		fmt.Println(tui.Red("Could not load TLS keys"))
		fmt.Println(tui.Red(err.Error()))
	}

	// Array of gRPC options with credentials
	serv.opts = []grpc.ServerOption{grpc.Creds(serv.creds),
		grpc.UnaryInterceptor(serv.UnaryInterceptor)}

	// Create the server object, attach all services
	serv.server = grpc.NewServer(serv.opts...)

	return nil
}

func (serv *ClientRPC) WriteConfig() error {

	var jsonData []byte
	jsonData, err := json.MarshalIndent(serv, "", "    ")
	if err != nil {
		fmt.Println("Error: Failed to write JSON data to server configuration file")
		fmt.Println(err)
	} else {
		_ = ioutil.WriteFile(serverFile, jsonData, 0755)
		fmt.Println(tui.Green("Server configuration file written."))
	}

	return nil
}

// ---------------------------------------------------
// AUTHENTICATION

func (serv *ClientRPC) AuthenticateClient(ctx context.Context, s *ClientRPC) (clearance string, admin bool, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		// Check for user registration, password and admin rights
		for _, u := range serv.Users {
			admin = u.Administrator
			if clientLogin == u.Name && clientPassword == u.PasswordHashString {
				clearance = "clear"
				return clearance, admin, nil
			}
			if clientLogin == u.Name && u.PasswordHashString == "" {
				fmt.Println("PasswordHash is empty")
				clearance = "reg"
				return clearance, admin, nil
			}
		}
	}
	clearance = "none"
	admin = false

	return clearance, admin, nil
}

func (serv *ClientRPC) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	clearance, admin, err := serv.AuthenticateClient(ctx, serv)
	if err != nil {
		return nil, err
	}

	// Save metadata to context for further processing
	ctx = context.WithValue(ctx, clearanceCtx, clearance)
	ctx = context.WithValue(ctx, adminCtx, admin)
	return handler(ctx, req)
}
