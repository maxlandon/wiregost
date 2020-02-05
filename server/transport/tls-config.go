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
	"crypto/x509"

	"github.com/maxlandon/wiregost/server/certs"
)

func getClientServerTLSConfig(host string) *tls.Config {

	caCertPtr, _, err := certs.GetCertificateAuthority(certs.UserCA)
	if err != nil {
		clientLog.Fatalf("Invalid ca type (%s): %v", certs.UserCA, host)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCertPtr)

	_, _, err = certs.UserServerGetCertificate(host)
	if err == certs.ErrCertDoesNotExist {
		certs.UserServerGenerateCertificate(host)
	}

	certPEM, keyPEM, err := certs.UserServerGetCertificate(host)
	if err != nil {
		clientLog.Errorf("Failed to generate or fetch certificate %s", err)
		return nil
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		clientLog.Fatalf("Error loading server certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:                  caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                caCertPool,
		Certificates:             []tls.Certificate{cert},
		CipherSuites:             []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}
