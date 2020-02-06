// Wiregost - Golang Exploitation Framework
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

package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	consts "github.com/maxlandon/wiregost/client/constants"
	ghostpb "github.com/maxlandon/wiregost/protobuf/ghost"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/log"
	"github.com/maxlandon/wiregost/server/rpc"
)

var (
	clientLog = log.ServerLogger("transport", "client")
)

var once = &sync.Once{}

// StartClientListener - Starts a MTLS listener, waiting for client connections
func StartClientListener(bindIface string, port uint16) (net.Listener, error) {

	clientLog.Infof("Starting Raw TCP/TLS listener on %s:%s", bindIface, port)

	tlsConfig := getClientServerTLSConfig(bindIface)
	ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", bindIface, port), tlsConfig)
	if err != nil {
		clientLog.Error(err)
		return nil, err
	}

	// go AcceptClientConnections(ln)
	return ln, nil
}

// AcceptClientConnections - Handle connections from clients
func AcceptClientConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errType, ok := err.(*net.OpError); ok && errType.Op == "accept" {
				break
			}
			clientLog.Errorf("Accept failed: %v", err)
			continue
		}
		go handleClientConnection(conn)
	}
}

// handleClientConnection - Authenticate connections, and bind events, commands, etc.
func handleClientConnection(conn net.Conn) {
	defer conn.Close()
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return
	}
	tlsConn.Read([]byte{}) // Unless you read 0 bytes the TLS handshake will not complete
	logState(tlsConn)
	certs := tlsConn.ConnectionState().PeerCertificates
	if len(certs) < 1 {
		return
	}
	user := certs[0].Subject.CommonName // Get user name from cert CN
	clientLog.Infof("Accepted incoming client connection: %s (%s)", conn.RemoteAddr(), user)

	log.AuditLogger.WithFields(logrus.Fields{
		"pkg":  "transport",
		"user": user,
	}).Info("connected")

	client := core.GetClient(certs[0])
	core.Clients.AddClient(client)

	core.EventBroker.Publish(core.Event{
		EventType: consts.JoinedEvent,
		Client:    client,
	})

	cleanup := func() {
		clientLog.Infof("Closing connection to client (%s)", client.User)
		log.AuditLogger.WithFields(logrus.Fields{
			"pkg":      "transport",
			"operator": client.User,
		}).Info("disconnected")
		core.Clients.RemoveClient(client.ID)
		conn.Close()
		core.EventBroker.Publish(core.Event{
			EventType: consts.LeftEvent,
			Client:    client,
		})
	}

	// Handle RPC requests/responses
	go func() {
		defer once.Do(cleanup)
		rpcHandlers := rpc.GetRPCHandlers()
		tunHandlers := rpc.GetTunnelHandlers()
		for {
			envelope, err := socketReadEnvelope(conn)
			if err != nil {
				clientLog.Errorf("Socket read error %v", err)
				return
			}
			// RPC
			if rpcHandler, ok := (*rpcHandlers)[envelope.Type]; ok {
				timeout := time.Duration(envelope.Timeout)
				go rpcHandler(envelope.Data, timeout, func(data []byte, err error) {
					errStr := ""
					if err != nil {
						errStr = fmt.Sprintf("%v", err)
					}
					client.Send <- &ghostpb.Envelope{
						ID:   envelope.ID,
						Data: data,
						Err:  errStr,
					}
				})
				log.AuditLogger.WithFields(logrus.Fields{
					"pkg":           "transport",
					"user":          client.User,
					"envelope_type": envelope.Type,
				}).Info("rpc command")
			} else if tunHandler, ok := (*tunHandlers)[envelope.Type]; ok {
				go tunHandler(client, envelope.Data, func(data []byte, err error) {
					errStr := ""
					if err != nil {
						errStr = fmt.Sprintf("%v", err)
					}
					client.Send <- &ghostpb.Envelope{
						ID:   envelope.ID,
						Data: data,
						Err:  errStr,
					}
				})
			} else {
				client.Send <- &ghostpb.Envelope{
					ID:   envelope.ID,
					Data: []byte{},
					Err:  "Unknown rpc command",
				}
			}
		}
	}()

	// Handle events send/receive
	events := core.EventBroker.Subscribe()
	defer core.EventBroker.Unsubscribe(events)
	go socketEventLoop(conn, events)

	defer once.Do(cleanup)
	for envelope := range client.Send {
		err := socketWriteEnvelope(conn, envelope)
		if err != nil {
			clientLog.Errorf("Socket error %v", err)
			return
		}
	}
}

func logState(tlsConn *tls.Conn) {
	clientLog.Debug(">>>>>>>>>>>>>>>> TLS State <<<<<<<<<<<<<<<<")
	state := tlsConn.ConnectionState()
	clientLog.Debugf("Version: %x", state.Version)
	clientLog.Debugf("HandshakeComplete: %t", state.HandshakeComplete)
	clientLog.Debugf("DidResume: %t", state.DidResume)
	clientLog.Debugf("CipherSuite: %x", state.CipherSuite)
	clientLog.Debugf("NegotiatedProtocol: %s", state.NegotiatedProtocol)
	clientLog.Debugf("NegotiatedProtocolIsMutual: %t", state.NegotiatedProtocolIsMutual)
	clientLog.Debug("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		clientLog.Debugf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		clientLog.Debugf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	clientLog.Debug(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}
