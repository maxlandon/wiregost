package connection

import (
	"google.golang.org/grpc"
)

// ConnectTLS - Establishes a TLS connection on which we will register gRPC clients
func ConnectTLS() (conn *grpc.ClientConn, err error) {

	// Load certificates required by Wiregost server

	// Dial server with these certificates
	conn, err = grpc.Dial("", grpc.WithInsecure())

	return
}
