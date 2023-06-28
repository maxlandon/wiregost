package multiplayer

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
	"encoding/json"
	"encoding/pem"
	"testing"

	clienttransport "github.com/maxlandon/wiregost/internal/client/transport"
	"github.com/maxlandon/wiregost/internal/server/certs"
)

func TestRootOnlyVerifyCertificate(t *testing.T) {
	certs.SetupCAs()

	data, err := NewOperatorConfig("zerocool", "localhost", uint16(1337))
	if err != nil {
		t.Fatalf("failed to generate test player profile %s", err)
	}
	config := &ClientConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		t.Fatalf("failed to parse client config %s", err)
	}

	_, _, err = certs.OperatorServerGetCertificate("localhost")
	if err == certs.ErrCertDoesNotExist {
		certs.OperatorServerGenerateCertificate("localhost")
	}

	// Test with a valid certificate
	certPEM, _, _ := certs.OperatorServerGetCertificate("localhost")
	block, _ := pem.Decode(certPEM)
	err = clienttransport.RootOnlyVerifyCertificate(config.CACertificate, [][]byte{block.Bytes})
	if err != nil {
		t.Fatalf("root only verify certificate error: %s", err)
	}

	// Test with wrong CA
	wrongCert, _ := certs.GenerateECCCertificate("https-ca", "foobar", false, false)
	block, _ = pem.Decode(wrongCert)
	err = clienttransport.RootOnlyVerifyCertificate(config.CACertificate, [][]byte{block.Bytes})
	if err == nil {
		t.Fatal("root only verify cert verified a certificate with invalid ca!")
	}
}
