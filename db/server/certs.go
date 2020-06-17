package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	serverpb "github.com/maxlandon/wiregost/proto/v1/gen/go/server"
)

type certServer struct {
	*serverpb.UnimplementedCertificateRPCServer
}

func (c *certServer) GetCertificate(context.Context, *serverpb.Get) (*serverpb.CertificateKeyPair, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCertificate not implemented")
}
func (c *certServer) RemoveCertificate(context.Context, *serverpb.Remove) (*serverpb.Removed, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveCertificate not implemented")
}
