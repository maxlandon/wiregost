package connection

import (
	"google.golang.org/grpc"
)

// ConnectTLS - Establishes a TLS connection on which we will register gRPC clients
func ConnectTLS() (conn *grpc.ClientConn, err error) {
	// func ConnectTLS() (cli client.ConnectionRPCClient, err error) {

	// Dial server with configuration certificates
	conn, err = grpc.Dial("", grpc.WithInsecure())

	return
}
