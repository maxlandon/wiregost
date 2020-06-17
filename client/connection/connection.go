package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maxlandon/wiregost/client/assets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024

	// ClientMaxReceiveMessageSize - Max gRPC message size ~2Gb
	ClientMaxReceiveMessageSize = 2 * gb

	defaultTimeout = time.Duration(10 * time.Second)
)

// ConnectTLS - Establishes a TLS connection on which we will register gRPC clients
func ConnectTLS() (conn *grpc.ClientConn, err error) {

	conf := assets.ServerConfig

	// Setup TLS
	tlsConfig, err := getTLSConfig(conf.CACertificate, conf.Certificate, conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Set gRPC options
	creds := credentials.NewTLS(tlsConfig)
	options := []grpc.DialOption{
		grpc.WithTimeout(defaultTimeout),
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	}

	// Dial server with these certificates
	server := fmt.Sprintf("%s:%d", conf.LHost, conf.LPort)
	conn, err = grpc.Dial(server, options...)

	fmt.Println("test")
	return
}

func getTLSConfig(caCertificate string, certificate string, privateKey string) (*tls.Config, error) {

	certPEM, err := tls.X509KeyPair([]byte(certificate), []byte(privateKey))
	if err != nil {
		log.Printf("Cannot parse client certificate: %v", err)
		return nil, err
	}

	// Load CA cert
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCertificate))

	// Setup config with custom certificate validation routine
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certPEM},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Don't worry I sorta know what I'm doing
		VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			return rootOnlyVerifyCertificate(caCertificate, rawCerts)
		},
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

// rootOnlyVerifyCertificate - Go doesn't provide a method for only skipping hostname validation so
// we have to disable all of the fucking certificate validation and re-implement everything.
// https://github.com/golang/go/issues/21971
func rootOnlyVerifyCertificate(caCertificate string, rawCerts [][]byte) error {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(caCertificate))
	if !ok {
		log.Printf("Failed to parse root certificate")
		os.Exit(3)
	}

	cert, err := x509.ParseCertificate(rawCerts[0]) // We should only get one cert
	if err != nil {
		log.Printf("Failed to parse certificate: " + err.Error())
		return err
	}

	// Basically we only care if the certificate was signed by our authority
	// Go selects sensible defaults for time and EKU, basically we're only
	// skipping the hostname check, I think?
	options := x509.VerifyOptions{
		Roots: roots,
	}
	if _, err := cert.Verify(options); err != nil {
		log.Printf("Failed to verify certificate: " + err.Error())
		return err
	}

	return nil
}
