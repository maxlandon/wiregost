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

package rpc

import (
	"bytes"
	"crypto/x509"
	"time"

	"github.com/gogo/protobuf/proto"

	clientpb "github.com/maxlandon/wiregost/protobuf/client"
	"github.com/maxlandon/wiregost/server/certs"
	"github.com/maxlandon/wiregost/server/core"
	"github.com/maxlandon/wiregost/server/users"
)

func rpcListUsers(data []byte, timeout time.Duration, resp RPCResponse) {
	clientCerts := certs.UserClientListCertificates()

	players := &clientpb.Players{Players: []*clientpb.Player{}}
	for _, cert := range clientCerts {
		players.Players = append(players.Players, &clientpb.Player{
			Client: &clientpb.Client{
				User: cert.Subject.CommonName,
			},
			Online: isPlayerOnline(cert),
		})
	}

	data, err := proto.Marshal(players)
	if err != nil {
		rpcLog.Errorf("Error encoding rpc response %v", err)
	}
	resp(data, err)
}

func rpcAddUser(data []byte, timeout time.Duration, resp RPCResponse) {

	userReq := &clientpb.UserReq{}
	err := proto.Unmarshal(data, userReq)
	if err != nil {
		resp(data, err)
	}

	err = users.NewUser(userReq.User, userReq.LHost, userReq.LPort, userReq.Default)

	userRes := &clientpb.User{}
	if err != nil {
		userRes.Success = false
		userRes.Err = err.Error()
	} else {
		userRes.Success = true
	}

	data, err = proto.Marshal(userRes)
	resp(data, err)
}

// isPlayerOnline - Is a player connected using a given certificate
func isPlayerOnline(cert *x509.Certificate) bool {
	for _, client := range *core.Clients.Connections {
		if client.Certificate == nil {
			continue // Server certificate is nil
		}
		if bytes.Equal(cert.Raw, client.Certificate.Raw) {
			return true
		}
	}
	return false
}
