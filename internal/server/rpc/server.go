package rpc

// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"context"
	"errors"
	"runtime"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/maxlandon/wiregost/internal/proto/clientpb"
	"github.com/maxlandon/wiregost/internal/proto/commonpb"
	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
	"github.com/maxlandon/wiregost/internal/server/log"
	"github.com/maxlandon/wiregost/internal/server/version"
)

var rpcLog = log.NamedLogger("rpc", "server")

const (
	minTimeout = time.Duration(30 * time.Second)
)

// Server - gRPC server
type Server struct {
	// Magical methods to break backwards compatibility
	// Here be dragons: https://github.com/grpc/grpc-go/issues/3794
	rpcpb.UnimplementedCoreServer
}

// NewServer - Create new server instance
func NewServer() *Server {
	// core.StartEventAutomation()
	return &Server{}
}

// GenericRequest - Generic request interface to use with generic handlers
type GenericRequest interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message

	GetRequest() *commonpb.Request
}

// GenericResponse - Generic response interface to use with generic handlers
type GenericResponse interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message

	GetResponse() *commonpb.Response
}

// GetVersion - Get the server version
func (rpc *Server) GetVersion(ctx context.Context, _ *commonpb.Empty) (*clientpb.Version, error) {
	dirty := version.GitDirty != ""
	semVer := version.SemanticVersion()
	compiled, _ := version.Compiled()
	return &clientpb.Version{
		Major:      int32(semVer[0]),
		Minor:      int32(semVer[1]),
		Patch:      int32(semVer[2]),
		Commit:     version.GitCommit,
		Dirty:      dirty,
		CompiledAt: compiled.Unix(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}, nil
}

func (rpc *Server) getClientCommonName(ctx context.Context) string {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	tlsAuth, ok := client.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ""
	}
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return ""
	}
	if tlsAuth.State.VerifiedChains[0][0].Subject.CommonName != "" {
		return tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	}
	return ""
}

// getTimeout - Get the specified timeout from the request or the default
func (rpc *Server) getTimeout(req GenericRequest) time.Duration {
	timeout := req.GetRequest().Timeout
	if time.Duration(timeout) < time.Second {
		return minTimeout
	}
	return time.Duration(timeout)
}

// getError - Check an implant's response for Err and convert it to an `error` type
func (rpc *Server) getError(resp GenericResponse) error {
	respHeader := resp.GetResponse()
	if respHeader != nil && respHeader.Err != "" {
		return errors.New(respHeader.Err)
	}
	return nil
}
