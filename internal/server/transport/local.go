package transport

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
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/maxlandon/wiregost/internal/proto/rpcpb"
	"github.com/maxlandon/wiregost/internal/server/log"
	"github.com/maxlandon/wiregost/internal/server/rpc"
)

const bufSize = 2 * mb

var bufConnLog = log.NamedLogger("transport", "local")

// LocalListener - Bind gRPC server to an in-memory listener, which is
//
//	typically used for unit testing, but ... it should be fine
func LocalListener() (*grpc.Server, *bufconn.Listener, error) {
	bufConnLog.Infof("Binding gRPC to listener ...")
	ln := bufconn.Listen(bufSize)
	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(ServerMaxMessageSize),
		grpc.MaxSendMsgSize(ServerMaxMessageSize),
	}
	options = append(options, initMiddleware(false)...)
	grpcServer := grpc.NewServer(options...)
	rpcpb.RegisterCoreServer(grpcServer, rpc.NewServer())
	go func() {
		panicked := true
		defer func() {
			if panicked {
				bufConnLog.Errorf("stacktrace from panic: %s", string(debug.Stack()))
			}
		}()
		if err := grpcServer.Serve(ln); err != nil {
			bufConnLog.Fatalf("gRPC local listener error: %v", err)
		} else {
			panicked = false
		}
	}()
	return grpcServer, ln, nil
}