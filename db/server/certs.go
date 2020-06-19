package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

type certServer struct {
	*serverpb.UnimplementedCertificateRPCServer
}

func (c *certServer) GetCertificate(ctx context.Context, in *serverpb.Get) (out *serverpb.CertificateKeyPair, err error) {

	// Generally the in certificate just contains a host name
	DB.Find(out).Where("hostname = ?", in.Cert.Hostname).Where("keytype = ?", in.Cert.KeyType)

	if out == nil {
		return nil, errors.New("No CertificateKeyPair matches for the certificate parameters given")
	}

	return out, nil
}

func (c *certServer) AddCertificate(ctx context.Context, in *serverpb.Add) (out *serverpb.Added, err error) {

	// Create CertificateKeyPair object
	cert := &serverpb.CertificateKeyPair{
		Hostname:    in.Hostname,
		KeyType:     in.KeyType,
		Certificate: in.Certificate,
		PrivateKey:  in.PrivateKey,
	}

	// Add it to DB
	errDB := DB.Create(cert)
	if len(errDB.GetErrors()) != 0 {
		return nil, errDB.GetErrors()[0]
	}

	out = &serverpb.Added{Cert: cert}
	return
}

func (c *certServer) RemoveCertificate(ctx context.Context, in *serverpb.Remove) (*serverpb.Removed, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveCertificate not implemented")
}
