package main

import (
	// Standard
	"crypto/tls"
	"fmt"
	"log"
	"os"

	// Wiregost
	"github.com/maxlandon/wiregost/internal/wiregost"
)

func main() {
	fmt.Println("Spectre server listening for connections")

	// TLS Config
	cert, err := tls.LoadX509KeyPair("/home/para/.wiregost/server/certificates/spectre_public.pem",
		"/home/para/.wiregost/server/certificates/spectre_private.pem")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", "localhost:7777", cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Change directory to Wiregost root project
	os.Chdir("/home/para/pentest/wiregost")

	// Start server
	wg := wiregost.NewWiregost()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		// server.Connect(conn)
		wg.Endpoint.Connect(conn)
	}
}
