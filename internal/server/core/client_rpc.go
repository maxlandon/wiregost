package core

// This file contains the code for the RPC facing the WireGost client.

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/evilsocket/islazy/tui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var userHomeDir, _ = os.UserHomeDir()

type ClientRPC struct {
	Protocol string
	Port     int
	CertPath string
	KeyPath  string
}

func NewClientRPC() *ClientRPC {
	serv := &ClientRPC{
		// Change this with parameters available in server-specific config file.
		Protocol: "tcp",
		Port:     7777,
		CertPath: userHomeDir + "/.wiregost/server/certificates/spectre.crt",
		KeyPath:  userHomeDir + "/.wiregost/server/certificates/spectre.key",
	}

	lis, err := net.Listen(serv.Protocol, fmt.Sprintf("%s:%d", "localhost", serv.Port))
	if err != nil {
		log.Fatalf("%s Failed to listen on port %d: %v", tui.RED, serv.Port, err)
	}

	// Register all ServiceManagers instances in WireGost
	userManager := UserManager{}

	// Create TLS Credentials
	creds, err := credentials.NewServerTLSFromFile(serv.CertPath, serv.KeyPath)
	if err != nil {
		fmt.Println(tui.Red("Could not load TLS keys"))
		fmt.Println(tui.Red(err.Error()))
	}

	// Array of gRPC options with credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// Create the server object, attach all services
	grpcServer := grpc.NewServer(opts...)

	// Attach all Services to the ClientRPC server
	RegisterUserManagerServer(grpcServer, &userManager)

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}

	return serv
}

// ---------------------------------------------------
// AUTHENTICATION
